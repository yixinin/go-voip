import pyaudio
import wave
import numpy as np
from socket import *

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

    h1 = frameType.to_bytes(length=1, byteorder='big', signed=True)
    h2 = dataType.to_bytes(length=1, byteorder='big', signed=True)
    head = bytes(h1+h2)
    print("head: ", head)
    ts = 0
    while 1:
        body = stream.read(BUFSIZE)
        # playStream.write(body)
        # continue
        timeStamp = ts.to_bytes(length=8, byteorder="big", signed=False)
        data = bytes(head+timeStamp+body)
       # print("send before", data.__len__())
        # n = tcpClient.send(data)
        tcpClient.sendall(data)
        # print(n)
        ts += 1
        # print("header:", data[:8+2])
        # print(body.__len__(), data.__len__())
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


if __name__ == '__main__':
    # 发送登录信息
    h1 = frameType.to_bytes(length=1, byteorder='big', signed=True)
    h2 = dataType.to_bytes(length=1, byteorder='big', signed=True)
    header = bytes()
    token = b'00000000000000000000000000000000'
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

    Capture(tcpClient)

    tcpClient.close()
