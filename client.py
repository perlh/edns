import threading
import requests
import socket


import random

def generate_random_string(length: int, chars: str) -> str:
    """
    Generate a random string with the given length and character set.
    """
    # Create a list to store the random characters
    random_chars = []

    # Generate a list of random characters
    for _ in range(length):
        random_chars.append(random.choice(chars))

    # Return the list of characters as a string
    return ''.join(random_chars)

def generate_user_string(user_id: str, length: int, chars: str) -> str:
    """
    Generate a random string for the given user.
    """
    # Create a string that combines the user ID and the random characters
    user_string = user_id + generate_random_string(length, chars)

    # Return the user string
    return user_string


def get_ip():
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    try:
        # 不需要连接到远程主机，仅仅发送数据报
        s.connect(('10.255.255.255', 1))
        IP = s.getsockname()[0]
    except:
        IP = '127.0.0.1'
    finally:
        s.close()
    return IP



import hashlib

def get_md5_hash(data: str) -> str:
    """
    Generate the MD5 hash of a string.
    """
    # Create an MD5 hash object
    m = hashlib.md5()

    # Use the hash object to generate the MD5 hash
    m.update(data.encode('utf-8'))

    # Return the hexadecimal representation of the hash
    return m.hexdigest()


def base58_decode(cipher_input: str) -> str:
    """
    base58编码典型应用是比特币钱包，与base64相比，去除了0、I、O、l、/ +等不易辨认的6个字符
    :param cipher_input: 输入base58编码值
    :return: base58的解码值
    @author hongfeiyinxue 2022-04-30-1651329023.0065577
    """
    #  定义base58的58个编码
    base58 = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
    cipher = cipher_input
    #  检查密文字符的有效性，密文字符必须是base58中的字符，否则返回提示
    bad = ''
    for item in cipher:
        if base58.find(item) == -1:
            bad += item
    if bad != '':
        return '不是有效的Base58编码，请仔留意如下字符：' + bad

    #    获取密文每个字符在base58中的下标值
    tmp = []
    for item in cipher:
        tmp.append(base58.find(item))
    temp = tmp[0]
    for i in range(len(tmp) - 1):
        temp = temp * 58 + tmp[i + 1]
    temp = bin(temp).replace('0b', '')
    #    print('temp=', temp, '-len-', len(temp))

    #   判断temp二进制编码数量能否被8整除，例如编码长度为18，首先截取（18%8）余数个字符求对应的ascii字符
    remainder = len(temp) % 8
    plain_text = ''

    if remainder != 0:
        temp_start = temp[0:remainder]
        plain_text = chr(int(temp[0:remainder], 2))

    for i in range(remainder, len(temp), 8):
        #    print(chr(int((temp[i:i+8]), 2)))
        plain_text += chr(int((temp[i:i + 8]), 2))
        i += 8
    return plain_text

def base58_encode(string_input):
    """
    base58编码典型应用是比特币钱包，与base64相比，去除了0、I、O、l、/ +等不易辨认的6个字符
    base58的编码思路是反复除以58取余数直至为0，base64的编码原理是64进制，2的6次方刚好等于64
    :param string_input: 输入待编码的字符
    :return: base58的编码值
    """
    base58 = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
    string = string_input
    string_binary = ''
    #   获取输入字符ascii码值的二进制码字符串，8个bit为一组
    for i in range(len(string)):
        string_binary = string_binary + str('{:0>8}'.format(bin(ord(string[i])).replace('0b', '')))
    string_decimal = int(string_binary, 2)
    string_58_list = []
    while True:
        string_58_list.insert(0, string_decimal % 58)
        string_decimal = string_decimal // 58
        if string_decimal == 0:
            break
    string_58 = ""
    for i in string_58_list:
        string_58 += base58[i]
    return string_58

# print(base58_decode(base58_encode('hongfei')))
char_array = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

def genToken(user,passwd):
    md5_encode = get_md5_hash(user+passwd)
    # print(md5_encode)
    tmp = user+":"+md5_encode

    return base58_encode(tmp)
    # return "dsafasd"

# 设置刷新时间间隔（5秒）
interval =20
# global old_ip

old_ip = ""


def refresh():
    global old_ip
    # 获取本机IP
    ip = get_ip()
    if old_ip != ip :
        old_ip = ip
        # 请求的URL
        url = "http://ddns.cool:10053/register_dns"
        email = "xxx@ddns.cool"
        passwd = "xxx"
        token = genToken(email, passwd)
        # print(token)
        # 打印IP地址
        # print(ip)
        # random_subdomain = generate_user_string("",10,char_array)
        # 请求的参数
        # print(random_subdomain)
        payload = {"token": token,"host_recode" : "zhoufang","recode_type":"a","recode_value":ip,"ttl":"600","domain":"ddns.hsm.cool"}
        # 设置定时器，定时调用refresh函数
        send_post(url,payload)
    threading.Timer(interval, refresh).start()


def send_post(url,payload):
    try:
        response = requests.post(url, data=payload,timeout=2)
        # response = requests.post(url, data=data, )
    except Exception as e:
        print(e)
        # 超时处理
        return ""


    # 发送POST请求

    print(response.text)
    # 打印响应内容
    return response.text

if __name__ == '__main__':
    # ip = get_ip()
    # print(ip)

    refresh()

