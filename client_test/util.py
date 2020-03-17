

def get_body_length(header):
    return int.from_bytes(header[2:], byteorder="little", signed=False)


def get_data_type(header):
    return header[1]


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
                leftBody = body[length:]
                if leftBody.__len__() > const.HEADER_LENGTH:  # 不止包含头部
                    header = leftBody[:const.HEADER_LENGTH]
                    # length = util.get_body_length()

                    # 如果不是一个完整包 合并到下一次
                    if leftBody.__len__() < const.HEADER_LENGTH + util.get_body_length(header):
                        preBuf = leftBody
                        read = 0  # 跳出while循环
                    else:
                        length = util.get_body_length(header)
                        print(header, length, "2----")
                        body = leftBody[const.HEADER_LENGTH:]
                        read = body.__len__()
                else:
                    preBuf = leftBody
                    read = 0  # 跳出while循环
    print("end recv buf")
