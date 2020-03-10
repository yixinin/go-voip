import cv2
import argparse
import os
import pdb
from socket import *


def Capture(tcpClient):
    cap = cv2.VideoCapture(0)
    while True:
        ret, frame = cap.read()
        cv2.imshow("video", frame)
        c = cv2.waitKey(1)
        if c == 27:
            break
    cap.release()
    cv2.destroyAllWindows()


def conn():
    # 发送登录信息
    header = bytes()
    uid = 1024
    cuid = 2048
    x1 = uid.to_bytes(length=8, byteorder='big', signed=True)
    x2 = cuid.to_bytes(length=8, byteorder='big', signed=True)
    header = bytes(x1+x2)
    HOST = '127.0.0.1'
    PORT = 9999
    BUFSIZE = 4096
    ADDR = (HOST, PORT)

    tcpClient = socket(AF_INET, SOCK_STREAM)
    tcpClient.connect(ADDR)

    tcpClient.send(header)
    Capture(tcpClient)

    tcpClient.close()


if __name__ == '__main__':
    Capture(0)
