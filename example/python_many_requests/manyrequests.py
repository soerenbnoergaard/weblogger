import http.client
import threading

def post(n):
    conn = http.client.HTTPConnection("localhost", 5000)
    conn.request("POST", "/?token=test&data={:d}".format(n), "")
    conn.close()

def get():
    conn = http.client.HTTPConnection("localhost", 5000)
    conn.request("GET", "/")
    conn.close()

post_threads = [threading.Thread(target=post, args=[n]) for n in range(10000)]
get_threads = [threading.Thread(target=get) for n in range(10000)]
for t1, t2 in zip(post_threads, get_threads):
    t1.start()
    t2.start()

for t1, t2 in zip(post_threads, get_threads):
    t1.join()
    t2.join()