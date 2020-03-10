from socket import *
import pyaudio
import wave
import numpy as np

frameType = 2  # ws frame type 1=text 2=binary
dataType = 1   # live data type 1=audio 2=video
BUFSIZE = 2048+2+8


def Recv(tcpClient):
    CHUNK = 512*2
    FORMAT = pyaudio.paInt16
    CHANNELS = 2
    RATE = 48000

    p = pyaudio.PyAudio()
    stream = p.open(rate=RATE, channels=CHANNELS, format=FORMAT, output=True)

    data = bytes()
    while 1:
        buf = bytes(data+tcpClient.recv(4096))
        siz = buf.__len__()
        while siz >= BUFSIZE:
            nbuf = buf[BUFSIZE:]
            stream.write(buf[8+2:BUFSIZE])
            buf = nbuf
            siz = nbuf.__len__()
        data = buf

    # 停止数据流
    stream.stop_stream()
    stream.close()

    # 关闭 PyAudio
    p.terminate()


if __name__ == "__main__":
    # 发送登录信息
    header = bytes()
    h1 = frameType.to_bytes(length=1, byteorder='big', signed=True)
    h2 = dataType.to_bytes(length=1, byteorder='big', signed=True)

    token = b'00000000000000000000000000000001'
    rid = 10240
    x2 = rid.to_bytes(length=8, byteorder='big', signed=True)
    header = bytes(h1+h2+token+x2)

    HOST = '127.0.0.1'
    PORT = 9901
    # BUFSIZE = 4096
    ADDR = (HOST, PORT)

    tcpClient = socket(AF_INET, SOCK_STREAM)
    tcpClient.connect(ADDR)
    tcpClient.send(header)

    Recv(tcpClient)

    tcpClient.close()
