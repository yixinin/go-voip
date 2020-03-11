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


def handle_tcp(tcp):
    p = pyaudio.PyAudio()
    stream = p.open(rate=const.RATE,
                    channels=const.CHANNELS,
                    format=const.FORMAT,
                    output=True)

    preBuf = bytes()
    print("start recv buf")
    while 1:
        # 读取header
        header = tcp.recv(const.HEADER_LENGTH)
        body_length = util.get_body_length(header)
        # 读取body
        body = bytes()
        read = 0
        while read < body_length:
            need_length = 0
            unread = body_length - read
            if unread > const.TCP_BUFSIZE:
                need_length = const.TCP_BUFSIZE
            else:
                need_length = unread

            sub_body = tcp.recv(need_length)
            body = (body+sub_body)
            read += sub_body.__len__()
        if header[1] == const.AUDIO_TYPE:
            print(body.__len__(), body_length)
            play_audio(stream, body)
        elif header[1] == const.VIDEO_TYPE:
            play_video(body)
            if cv2.waitKey(1) & 0xFF == ord('q'):
                break

    print("end recv buf")


def play_audio(stream, body):
    if body.__len__() == 0:
        return

    stream.write(body)


def play_video(body):
    arr = np.frombuffer(body, np.uint8)
    frame = cv2.imdecode(arr, cv2.IMREAD_COLOR)
    cv2.imshow("video", frame)


def conn(user):
 # 发送登录信息
    loginBuf = const.get_user_header(user)

    ADDR = (const.HOST, const.TCP_PORT)

    tcpClient = socket(AF_INET, SOCK_STREAM)
    tcpClient.connect(ADDR)

    tcpClient.send(loginBuf)
    return tcpClient


def main():
    user = 2
    tcpClient = conn(user)

    # th_handle = threading.Thread(target=handle_tcp, args=(tcpClient,))
    th_audio = threading.Thread(target=capture_audio, args=(tcpClient,))
    # th_video = threading.Thread(target=capture_video, args=(tcpClient,))

    # th_handle.start()
    th_audio.start()
    # th_video.start()

    # th_handle.join()

    handle_tcp(tcpClient)

    # th_video.join()
    for th in (th_audio,):
        th.join()

    tcpClient.close()


if __name__ == '__main__':
    main()
