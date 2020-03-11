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
    ts = 0
    while 1:
        body = stream.read(BUFSIZE)
        # playStream.write(body)
        # continue
        timeStamp = ts.to_bytes(length=8, byteorder="little", signed=False)
        length = body.__len__().to_bytes(length=2, byteorder="little", signed=False)
        data = bytes(head + length + timeStamp + body)

        tcpClient.sendall(data)

        ts += 1

    print("end")
    stream.stop_stream()
    stream.close()
    p.terminate()
    # wf = wave.open(WAVE_OUTPUT_FILENAME, 'wb')
    # wf.setnchannels(CHANNELS)
    # wf.setsampwidth(p.get_sample_size(FORMAT))
    # wf.setframerate(RATE)
    # wf.writeframes(b''.join(frames))
    # wf.close()


def Recv(tcpClient):
    CHUNK = 512*2
    FORMAT = pyaudio.paInt16
    CHANNELS = 1
    RATE = 48000

    p = pyaudio.PyAudio()
    stream = p.open(rate=RATE, channels=CHANNELS, format=FORMAT, output=True)

    data = bytes()
    while 1:
        buf = (data + tcpClient.recv(4096))
        BUFSIZE = 12 + \
            int.from_bytes(buf[2:4], byteorder="little", signed=False)
        siz = buf.__len__()
        while siz >= BUFSIZE:
            stream.write(buf[8+4:BUFSIZE])
            buf = buf[BUFSIZE:]
            siz = buf.__len__()
        data = buf

    # 停止数据流
    stream.stop_stream()
    stream.close()

    # 关闭 PyAudio
    p.terminate()


if __name__ == '__main__':
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
