import requests
import tkinter as tk
import json

class ValidatorApp:
    def __init__(self, blockchain_server):
        self.blockchain_server = blockchain_server
        self.root = tk.Tk()
        self.root.title("Aplicação de Validadores")

        self.blocks_text = tk.Text(self.root, height=20, width=50)
        self.blocks_text.grid(row=0, column=0, padx=10, pady=10)

        self.refresh_button = tk.Button(self.root, text="Atualizar", command=self.refresh_blocks)
        self.refresh_button.grid(row=1, column=0, padx=10, pady=10)

        self.root.mainloop()

    def refresh_blocks(self):
        response = requests.get(self.blockchain_server + '/chain')
        if response.status_code == 200:
            chain = response.json()['chain']
            self.blocks_text.delete('1.0', tk.END)
            self.blocks_text.insert(tk.END, json.dumps(chain, indent=2))
        else:
            self.blocks_text.delete('1.0', tk.END)
            self.blocks_text.insert(tk.END, "Erro ao obter a blockchain")


blockchain_server = 'http://0.0.0.0:5000'  


validator_app = ValidatorApp(blockchain_server)
