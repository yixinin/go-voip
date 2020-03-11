import cv2
import os
from socket import *
import numpy as np
import const
import threading
import pyaudio
import util


def capture_video(tcpClient):
    cap = cv2.VideoCapture(0)
    while True:
        ret, frame = cap.read()
        if not ret:
            print("no video device")
            return
        img = cv2.imencode(".jpg", frame)[1]
        body = np.array(img).tobytes()

        size = body.__len__()
        header = const.get_video_header(size)  # 1+1+4
        buf = (header+body)

        tcpClient.sendall(buf)
        # print("send video buf", buf.__len__(), "\n")

    cap.release()


def capture_audio(tcpClient):
    p = pyaudio.PyAudio()
    stream = p.open(format=const.FORMAT,
                    channels=const.CHANNELS,
                    rate=const.RATE,
                    input=True,
                    frames_per_buffer=const.CHUNK)

    while 1:
        body = stream.read(const.AUDIO_BUFSIZE)

        size = body.__len__()
        header = const.get_audio_header(size)
        buf = bytes(header + body)

        tcpClient.sendall(buf)
        # print("send audio buf", buf.__len__(), "\n")

    stream.stop_stream()
    stream.close()
    p.terminate()


def play_audio(stream, body):
    print(body.__len__())
    stream.write(body)


def play_video(body):
    arr = np.frombuffer(body, np.uint8)
    frame = cv2.imdecode(arr, cv2.IMREAD_COLOR)
    cv2.imshow("video", frame)


def handle_buffer(tcpClient):
    p = pyaudio.PyAudio()
    stream = p.open(rate=const.RATE,
                    channels=const.CHANNELS,
                    format=const.FORMAT,
                    output=True)

    preBuf = bytes()
    print("start recv buf")
    while 1:
        buf = (preBuf + tcpClient.recv(const.TCP_BUFSIZE))
        header = buf[:2+4]
        length = util.get_body_length(header)
        body = buf[2+4:]

        print(header, length, "1----")
        read = body.__len__()  # 已读取的长度

        if read == length:
            if header[1] == const.AUDIO_TYPE:
                play_audio(stream, body)
                print("play 1")
            elif header[1] == const.VIDEO_TYPE:
                play_video(body)
                if cv2.waitKey(1) & 0xFF == ord('q'):
                    break

        elif read < length:  # 报文太短 拆包
            # 拆包 合并
            preBuf = buf  # 下次再读取

        elif read > length:  # 报文太长 粘包
            # 粘包  分解

            while read > length:

                # 先读取第一个包
                if header[1] == const.AUDIO_TYPE:
                    play_audio(stream, body[:length])
                    print("play 3")
                elif header[1] == const.VIDEO_TYPE:
                    play_video(body[:length])
                    if cv2.waitKey(1) & 0xFF == ord('q'):
                        break

                # 读取剩余
                if body.__len__() > length+2+4:  # 不止包含头部
                    header = body[length: length + 2+4]
                    # length = util.get_body_length()

                    # 如果不是一个完整包 合并到下一次
                    if body.__len__() < length + util.get_body_length() + 2+4:
                        preBuf = body[length:]
                        read = 0  # 跳出while循环
                    else:
                        length = util.get_body_length()
                        print(header, length, "2----")
                        body = body[length+2+4:]
                        read = body.__len__()
    print("end recv buf")


def conn(user):
 # 发送登录信息
    loginBuf = const.get_user_header(user)

    ADDR = (const.HOST, const.TCP_PORT)

    tcpClient = socket(AF_INET, SOCK_STREAM)
    tcpClient.connect(ADDR)

    tcpClient.send(loginBuf)
    return tcpClient


if __name__ == '__main__':
    user = 1
    tcpClient = conn(user)

    # th_handle = threading.Thread(target=handle_buffer, args=(tcpClient,))
    th_audio = threading.Thread(target=capture_audio, args=(tcpClient,))
    th_video = threading.Thread(target=capture_video, args=(tcpClient,))

    # th_handle.start()
    # th_audio.start()
    # th_video.start()

    # th_handle.join()

    handle_buffer(tcpClient)

    # th_video.join()
    th_audio.join()

    tcpClient.close()
