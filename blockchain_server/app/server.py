from flask import Flask, jsonify, request
from blockchain import Blockchain
import os

app = Flask(__name__)
blockchain = Blockchain()


from flask_cors import CORS
CORS(app)


node_identifier = str(hash(os.environ.get('NODE_IDENTIFIER', 'default_identifier')))

@app.route('/blocks', methods=['GET'])
def get_blocks():
    response = {
        'blocks': blockchain.get_blocks(),
        'length': len(blockchain.get_blocks())
    }
    return jsonify(response), 200
    
@app.route('/mine', methods=['POST'])
def mine():
    data = request.get_json()

    required_fields = ['voter_id', 'candidate_id']
    if not all(field in data for field in required_fields):
        return jsonify({'error': 'Campos voter_id e candidate_id são obrigatórios'}), 400

    voter_id = data['voter_id']
    candidate_id = data['candidate_id']

    last_block = blockchain.last_block
    last_proof = last_block['proof']
    proof = blockchain.proof_of_work(last_proof)

    blockchain.new_transaction(
        voter_id=voter_id,
        candidate_id=candidate_id,
    )

    previous_hash = blockchain.hash(last_block)
    block = blockchain.new_block(proof, previous_hash)

    response = {
        'message': "Novo bloco criado",
        'index': block['index'],
        'transactions': block['transactions'],
        'proof': block['proof'],
        'previous_hash': block['previous_hash'],
    }
    return jsonify(response), 200

@app.route('/chain', methods=['GET'])
def full_chain():
    response = {
        'chain': blockchain.chain,
        'length': len(blockchain.chain),
    }
    return jsonify(response), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
