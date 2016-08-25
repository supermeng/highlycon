import gevent
import time
import requests

from gevent import monkey
monkey.patch_socket()

url = "http://127.0.0.1:8888"


def sigle_test():
    session = requests.Session()
    for i in range(1000):
        session.get(url)
        gevent.sleep(0)


def con_test():
    start = time.time()
    gs = []
    for i in range(100):
        gs.append(gevent.spawn(sigle_test))
    gevent.joinall(gs)
    duration = time.time() - start
    print 'duration:', duration
    print 'qps:', (100000 / duration)

con_test()
