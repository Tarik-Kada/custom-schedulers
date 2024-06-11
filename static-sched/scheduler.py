from waitress import serve
from flask import Flask, jsonify, request
import sys
import time

app = Flask(__name__)

@app.route('/', methods=['POST'])
def podScheduler():
    print("---")
    # Start timer to measure total time to get to scheduling decision
    start = time.time()
    # Read the request body
    request_data = request.get_json()

    # Get all node names in the cluster:
    nodeInfo = request_data['clusterInfo']['nodes']
    nodes = [node['nodeName'] for node in nodeInfo if node['nodeStatus'] == 'Ready']
    # Read the Paramaters field from the request body
    param = request_data['parameters']

    # Get the worker number from the request body and map it on so that 3 goes to 3 and 4 goes to 1
    nodes = sorted(nodes)
    selected_node = nodes[param['worker'] % len(nodes)]

    # Print the total time to get to scheduling decision
    print(f"Time to get to scheduling decision: {time.time() - start}", file=sys.stdout)
    sys.stdout.flush()
    return jsonify({"node": selected_node})

if __name__ == '__main__':
    print("Static scheduler acitve. Listening on port: 5000", file=sys.stdout)
    sys.stdout.flush()
    serve(app, host='0.0.0.0', port=5000)
