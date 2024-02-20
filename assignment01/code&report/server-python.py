###############################################################################
# server-python.py
# Name:
# NetId:
###############################################################################

from pickle import NONE
import sys
import socket

RECV_BUFFER_SIZE = 2048
QUEUE_LENGTH = 10

def server(server_port):
    """TODO: Listen on socket and print received message to sys.stdout"""
    l = None
    try:
        l = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        l.bind(('', int(server_port)))
        l.listen(20)
    except socket.error as msg:
        if l:
            l.close()
        print('could not open socket:', msg)
        sys.exit(1)

    while True:
        conn, _ = l.accept()
        n = 1
        while n > 0:
            data = conn.recv(2048)
            if not data:
                break
            sys.stdout.buffer.write(data)
            sys.stdout.buffer.flush()

    conn.close()

    l.close()

    


def main():
    """Parse command-line argument and call server function """
    if len(sys.argv) != 2:
        sys.exit("Usage: python server-python.py [Server Port]")
    server_port = int(sys.argv[1])
    server(server_port)

if __name__ == "__main__":
    main()



    