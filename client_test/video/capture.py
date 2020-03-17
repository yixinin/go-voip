import cv2
import argparse
import os
import pdb
from socket import *
import numpy as np
from PIL import Image

frameType = 2  # ws frame type 1=text 2=binary
dataType = 1   # live data type 1=audio 2=video


def Capture(tcpClient):

    h1 = frameType.to_bytes(length=1, byteorder='little', signed=True)
    h2 = dataType.to_bytes(length=1, byteorder='little', signed=True)
    head = bytes(h1+h2)

    cap = cv2.VideoCapture(0)

    while True:
        ret, frame = cap.read()
        if not ret:
            return
        # imgRGB = cv2.cvtColor(frame, cv2.IMREAD_COLOR)
        # r, buf = cv2.imencode(".jpg", imgRGB)
        # body = Image.fromarray(np.uint8(buf)).tobytes()
        img = cv2.imencode(".jpg", frame)[1]
        body = np.array(img).tobytes()

        # play(body)
        # if cv2.waitKey(1) & 0xFF == ord('q'):
        #     break
        # continue

        length = body.__len__().to_bytes(length=4, byteorder='little', signed=False)
        header = (head + length)  # 1+1+4
        buf = (header+body)

        tcpClient.sendall(buf)
        # continue

        # print(ts, header, buf.__len__())
        # siz = buf.__len__()
        # while siz > 0:
        #     if siz > 4096:
        #         n = tcpClient.send(buf[:4096])
        #         buf = buf[n:]
        #     else:
        #         n = tcpClient.send(buf)
        #         buf = buf[n:]
        #     siz = buf.__len__()

        # ts += 1
        # c = cv2.waitKey(1)
        # if c == 'c':
        #     break
    cap.release()


def play(data):
    arr = np.frombuffer(data, np.uint8)
    frame = cv2.imdecode(arr, cv2.IMREAD_COLOR)

    cv2.imshow("video", frame)


def conn():
 # 发送登录信息
    h1 = frameType.to_bytes(length=1, byteorder='little', signed=True)
    h2 = dataType.to_bytes(length=1, byteorder='little', signed=True)
    header = bytes()
    token = b'00000000000000000000000000000000'
    rid = 10240
    x2 = rid.to_bytes(length=8, byteorder='little', signed=True)
    header = bytes(h1+h2+token+x2)
    HOST = '127.0.0.1'
    PORT = 9901
    # BUFSIZE = 4096
    ADDR = (HOST, PORT)

    tcpClient = socket(AF_INET, SOCK_STREAM)
    tcpClient.connect(ADDR)

    tcpClient.send(header)
    return tcpClient


if __name__ == '__main__':
    tcpClient = conn()

    Capture(tcpClient)

    tcpClient.close()
