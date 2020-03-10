import pyaudio
import wave
import numpy as np
from socket import *

BUFSIZE = 4096+2
frameType = 2  # ws frame type 1=text 2=binary
dataType = 1   # live data type 1=audio 2=video


def Capture(tcpClient):
    CHUNK = 512*2
    FORMAT = pyaudio.paInt16
    CHANNELS = 2
    RATE = 48000
    RECORD_SECONDS = 5
    WAVE_OUTPUT_FILENAME = "cache.wav"
    p = pyaudio.PyAudio()
    stream = p.open(format=FORMAT,
                    channels=CHANNELS,
                    rate=RATE,
                    input=True,
                    frames_per_buffer=CHUNK)
    print("开始缓存录音")
    # frames = []

    h1 = frameType.to_bytes(length=1, byteorder='big', signed=True)
    h2 = dataType.to_bytes(length=1, byteorder='big', signed=True)

    while (True):
        print('begin ')
        # for i in range(0, 100):
        data = stream.read(BUFSIZE)
        # frames.append(data)
        tcpClient.send(bytes(h1+h2+data))

        # recvData = tcpClient.recv(4096)
        # print(recvData.__len__())
        # audio_data = np.fromstring(data, dtype=np.short)
        # large_sample_count = np.sum(audio_data > 800)
        # temp = np.max(audio_data)
        # if temp > 800:
        #     print("检测到信号")
        #     print('当前阈值：', temp)
        #     break
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
