from flask import Flask, jsonify, request
from multiprocessing import Value
counter = Value('i', 0)
app = Flask(__name__)
app.config["JSON_AS_ASCII"] = False
timer = {}
@app.route("/", methods=["GET","POST"])
def out():
    js = request.get_json()
    print(js)
    return jsonify({"hello": "python mockerserver is running!"})
@app.route("/task", methods=["GET","POST"])
def task():
    param = request.args
    appName = param["appinstance_name"]
    with counter.get_lock():
        counter.value += 1
        print(counter.value)
    if counter.value % 300  == 0:
        return """{"ret_code":0,"err_msg":"","data":{"compute_cluster_name":"Mesos_formal_ps2","framework_name":"PS_Lite_Test2","appinstance_id":"t_nni_kwinsheng_413687274_00000000","appinstance_name":"t_nni_kwinsheng_413687274_00000000","appinstance_status":"success","appinstance_progress":0,"appinstance_start_time":"2019-12-25 21:28:32","appinstance_stop_time":"0000-00-00 00:00:00","dispatcher_name":"","exit_code":0,"exit_code_info":"","exit_error_source":0,"user_task_ret_code":-9999}}"""
    else:
        return """{"ret_code":0,"err_msg":"","data":{"compute_cluster_name":"Mesos_formal_ps2","framework_name":"PS_Lite_Test2","appinstance_id":"t_nni_kwinsheng_413687274_00000000","appinstance_name":"t_nni_kwinsheng_413687274_00000000","appinstance_status":"running","appinstance_progress":0,"appinstance_start_time":"2019-12-25 21:28:32","appinstance_stop_time":"0000-00-00 00:00:00","dispatcher_name":"","exit_code":0,"exit_code_info":"","exit_error_source":0,"user_task_ret_code":-9999}}"""
@app.route("/content", methods=["GET","POST"])
def content():
    param = request.args
    offset = int(param["offset"])
    lenth = int(param["length"])
    with open("./stdout") as f:
        f.seek(offset)
        data = f.read(lenth)
    return jsonify({"ret_code": 0, "err_msg": "", "data": {"sLog": data}})
if __name__ == "__main__":
    app.run(host="0.0.0.0", port="12123", debug=True)