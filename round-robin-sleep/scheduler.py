from waitress import serve
from flask import Flask, jsonify, request
import sys
import time

app = Flask(__name__)

@app.route('/', methods=['POST'])
def podScheduler():
    global worker
    print("---")
    # Start timer to measure total time to get to scheduling decision
    start = time.time()
    # Read the request body
    request_data = request.get_json()

    # Read the nodes field from the request body
    nodeInfo = request_data['clusterInfo']['nodes']
    nodes = [node['nodeName'] for node in nodeInfo if node['nodeStatus'] == 'Ready']
    # Order the nodes based on the node name in alphabetical order
    nodes = sorted(nodes)
    selected_node = nodes[worker]
    print(f"Selected Node: {selected_node}", file=sys.stdout)
    # Increment the worker number by 1
    worker = (worker + 1) % len(nodes)

    param = request_data['parameters']
    if 'sleepTime' in param:
        sleepTime = int(param['sleepTime'])
        print(f"Sleeping for {sleepTime} seconds...", file=sys.stdout)
        time.sleep(sleepTime)

    # Print the total time to get to scheduling decision
    print(f"Time to get to scheduling decision: {time.time() - start}", file=sys.stdout)
    sys.stdout.flush()
    return jsonify({"node": selected_node})

if __name__ == '__main__':
    global worker
    worker = 0
    print("Round Robin scheduler acitve. Listening on port: 5000", file=sys.stdout)
    sys.stdout.flush()
    serve(app, host='0.0.0.0', port=5000)
