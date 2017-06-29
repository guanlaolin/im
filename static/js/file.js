//发送文件
function sendFile() {
	var step = 1024; //每次读取的片段大小
	var position = 0; //读取起始位置
	var loaded = 0; //已读取的段数
	var slice = null; //文件分片对象
	var seg = 0;	//当前读取片数
	var seg_total = 0;	//总片数

	//获取文件信息
	var efile = document.getElementById("upload-file");
	var file = efile.files[0];
	var size = file.size; //获取文件大小
	
	seg_total = parseInt(size / step);

	//循环处理文件
	while (position <= size) {
		var reader = new FileReader();
		
		slice = null;

		//处理最后一片
		
		if ((position + step) > size) {
			slice = file.slice(position, size);
		} else {
			slice = file.slice(position, position+step);
		}
		//移动到下一片
		position = position + step;

		//二进制读取文件
		reader.readAsBinaryString(slice);

		reader.onload = function (evt) {
			var binary = evt.target.result;
			var msg = {
				"type": MSG_FILE,
				"from": parseInt($.cookie('uid')),
				"to": fid,
				"value": window.btoa(binary),
				"fileinfo": {
					"name": file.name,
					"segment":seg++,
					"total":seg_total
				}
			};
			//发送文件
			ws.send(JSON.stringify(msg));
		}
	}
	alert("发送文件成功");
}
