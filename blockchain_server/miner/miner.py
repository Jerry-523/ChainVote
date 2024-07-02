import hashlib
import requests
import time

# Blockchain server URL
BLOCKCHAIN_SERVER_URL = "http://localhost:5000"

def get_last_block():
    response = requests.get(f"{BLOCKCHAIN_SERVER_URL}/blocks")
    if response.status_code == 200:
        blocks = response.json().get("blocks")
        if blocks:
            return blocks[-1]
    return None

def proof_of_work(last_proof):
    proof = 0
    while not valid_proof(last_proof, proof):
        proof += 1
    return proof

def valid_proof(last_proof, proof):
    guess = f"{last_proof}{proof}".encode()
    guess_hash = hashlib.sha256(guess).hexdigest()
    return guess_hash[:4] == "0000"

def mine_block(voter_id, candidate_id):
    last_block = get_last_block()
    if last_block is None:
        print("Error getting the last block.")
        return

    last_proof = last_block.get("proof")
    if last_proof is None:
        print("The last block does not contain the proof of work ('proof').")
        return

    proof = proof_of_work(last_proof)

    data = {
        "voter_id": voter_id,
        "candidate_id": candidate_id
    }

    response = requests.post(f"{BLOCKCHAIN_SERVER_URL}/mine", json=data)
    if response.status_code == 200:
        print("Block mined successfully:", response.json())
    else:
        print("Error sending the new block:", response.json())

if __name__ == "__main__":
    voter_id = "voter123"
    candidate_id = "candidate456"

    while True:
        mine_block(voter_id, candidate_id)
        time.sleep(10)  # 10 seconds interval between mining blocks
