import pyaudio
import wave
import numpy as np
from socket import *
import threading

BUFSIZE = 512
frameType = 2  # ws frame type 1=text 2=binary
dataType = 1   # live data type 1=audio 2=video


def Capture(tcpClient):
    CHUNK = 512*2
    FORMAT = pyaudio.paInt16
    CHANNELS = 1
    RATE = 48000
    # RECORD_SECONDS = 5
    # WAVE_OUTPUT_FILENAME = "cache.wav"
    p = pyaudio.PyAudio()
    stream = p.open(format=FORMAT,
                    channels=CHANNELS,
                    rate=RATE,
                    input=True,
                    frames_per_buffer=CHUNK)
    # playStream = p.open(rate=RATE, channels=CHANNELS,
    #                     format=FORMAT, output=True)
    print("开始缓存录音")
    # frames = []

    h1 = frameType.to_bytes(length=1, byteorder='little', signed=True)
    h2 = dataType.to_bytes(length=1, byteorder='little', signed=True)
    head = bytes(h1+h2)
    print("head: ", head)

    while 1:
        body = stream.read(BUFSIZE)

        length = body.__len__().to_bytes(length=4, byteorder="little", signed=False)
        data = bytes(head + length + body)

        tcpClient.sendall(data)

    print("end")
    stream.stop_stream()
    stream.close()
    p.terminate()


def Recv(tcpClient):
    CHUNK = 512*2
    FORMAT = pyaudio.paInt16
    CHANNELS = 1
    RATE = 48000

    p = pyaudio.PyAudio()
    stream = p.open(rate=RATE, channels=CHANNELS, format=FORMAT, output=True)

    preBuf = bytes()
    while 1:
        buf = (preBuf + tcpClient.recv(BUFSIZE))
        header = buf[:2+4]
        body = buf[2+4:]

        length = int.from_bytes(header[2:], byteorder="little", signed=False)
        read = body.__len__()  # 已读取的长度

        if read == length:
            play(body)

        elif read < length:
            # 拆包 合并
            while read < length:
                unRead = length - read
                subBody = tcpClient.recv(BUFSIZE)

                read += subBody.__len__()
                if read >= length:
                    body = (body+subBody[:length-read])
                    preBuf = subBody[length-read:]
                    read == length
                else:
                    body = (body+subBody)
            play(body)

        elif read > length:
            # 粘包  分解
            while read > length:
                play(body[:length])

                body = body[length:]
                read -= length
            preBuf = body

    # 停止数据流
    stream.stop_stream()
    stream.close()

    # 关闭 PyAudio
    p.terminate()


def play(stream, body):
    stream.write(body)


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

    # th_capture = threading.Thread(target=Capture, args=(tcpClient,))
    th_recv = threading.Thread(target=Recv, args=(tcpClient,))
    Capture(tcpClient)
    # 启动线程
    # th_capture.start()
    th_recv.start()

    # th_capture.join()
    th_recv.join()
    # Capture(tcpClient)

    tcpClient.close()
