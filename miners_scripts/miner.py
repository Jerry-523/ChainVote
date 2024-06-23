import hashlib
import requests
import json

def mine_block():
    try:
        print("Solicitando o último bloco da cadeia...")
        response = requests.get('http://localhost:5000/chain')
        response.raise_for_status()  
        print("Resposta recebida: ", response.text)
        last_block = response.json()['chain'][-1]
        last_proof = last_block['proof']
        proof = proof_of_work(last_proof)

        new_block = {
            'proof': proof,
            'previous_hash': hashlib.sha256(json.dumps(last_block, sort_keys=True).encode()).hexdigest(),
        }
        return new_block
    except requests.exceptions.RequestException as e:
        print(f"Erro na requisição: {e}")
        return None
    except (KeyError, IndexError, json.JSONDecodeError) as e:
        print(f"Erro ao processar a resposta JSON: {e}")
        return None

def proof_of_work(last_proof):
    proof = 0
    while not valid_proof(last_proof, proof):
        proof += 1
    return proof

def valid_proof(last_proof, proof):
    guess = f'{last_proof}{proof}'.encode()
    guess_hash = hashlib.sha256(guess).hexdigest()
    return guess_hash[:4] == "0000"

if __name__ == '__main__':
    new_block = mine_block()
    if new_block:
        print("Novo bloco minerado: ", new_block)
        try:
            headers = {'Content-Type': 'application/json'}
            data = {
                'voter_id': 'test_voter',
                'candidate_id': 'test_candidate',
                'proof': new_block['proof'],
                'previous_hash': new_block['previous_hash']
            }
            response = requests.post('http://localhost:5000/mine', headers=headers, data=json.dumps(data))
            print("Resposta do servidor de mineração: ", response.text)
            print(response.json())
        except requests.exceptions.RequestException as e:
            print(f"Erro na requisição: {e}")
        except json.JSONDecodeError as e:
            print(f"Erro ao decodificar JSON: {e}")
    else:
        print("Não foi possível minerar o bloco.")
