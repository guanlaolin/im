remotePC = new RTCPeerConnection();
pc = new RTCPeerConnection();

pc.onicecandidate = function (event) {
	if (event.candidate !== null) {
		ws.send(JSON.stringify({
				"type": MSG_VIDEO,
				"from": parseInt($.cookie('uid')),
				"to": fid,
				"value": {
					'type': 'candidate_offer',
					'value': event.candidate
				}
			}));
	}
};

remotePC.onicecandidate = function (event) {
	if (event.candidate !== null) {
		ws.send(JSON.stringify({
				"type": MSG_VIDEO,
				"from": parseInt($.cookie('uid')),
				"to": fid,
				"value": {
					'type': 'candidate_anwser',
					'value': event.candidate
				}
			}));
	}
};

//收到远程视频流
pc.onaddstream = function (event) {
	console.log('onaddstream');
	document.getElementById('remoteVideo').src = URL.createObjectURL(event.stream);
};
remotePC.onaddstream = function (event) {
	console.log('remote onaddstream');
	document.getElementById('remoteVideo').src = URL.createObjectURL(event.stream);
};

window.onload = function () {
	db = window.openDatabase("db_ws_im", "5.0", "db for im", 1024 * 1024);

	//创建数据表
	var sql = 'CREATE TABLE tb_history ("type" INTEGER, "from" INTEGER, "to" INTEGER, "value" TEXT, "flag" INTEGER)';
	db.transaction(function (tx) {
		tx.executeSql(sql);
	});
	
	//emoji
	$('#right-edit').emojiarea({button:'#emoji'});

	//视频聊天
	//getUserMedia = (navigator.getUserMedia || navigator.webkitGetUserMedia || navigator.mozGetUserMedia || navigator.msGetUserMedia);

	//ws连接
	ws = new WebSocket("wss://192.168.43.219:8888/ws");
	//ws = new WebSocket("wss://127.0.0.1:8888/ws");
	ws.onopen = function () {
		var str = {
			"type": MSG_IDENTITY,
			"from": parseInt($.cookie('uid'))
		};
		ws.send(JSON.stringify(str));
		var p = document.createElement('p');
		p.innerHTML = "连接服务器成功";
		document.getElementById('right-show').appendChild(p);
	}
	ws.onerror = function () {
		var p = document.createElement('p');
		p.innerHTML = "与服务器通信失败";
		document.getElementById('right-show').appendChild(p);
		alert("ws error");
	}
	ws.onclose = function () {
		var p = document.createElement('p');
		p.innerHTML = "从服务器断开";
		document.getElementById('right-show').appendChild(p);
	}
	ws.onmessage = function (evt) {
		document.getElementById("ring").play();
		
		//解析消息
		var obj = JSON.parse(evt.data);
		var type = obj.type;
		var vfid = obj.from;
		var value = obj.value;

		console.log(evt.data);

		switch (type) {
		case MSG_TEXT:
			//插入历史表
			addMsg(obj.type, obj.to, obj.from, value, MSG_RCV);
			//正打开对应id页面
			if (fid == vfid) {
				//直接显示
				queryHistory(obj.from);
			} else {
				//修改未读消息小红点
				var count = parseInt($('#unread_' + vfid.toString()).text());
				$('#unread_' + vfid.toString()).text((++count).toString());
				$('#unread_' + vfid.toString()).show();
			}
			break;
		case MSG_FILE:
			//插入历史表
			addMsg(obj.type, obj.to, obj.from, '我给你发送了文件，点击链接下载 <a href=' + value + ' download=' + obj.fileinfo.name + '>' + obj.fileinfo.name + '</a>', MSG_RCV);
			//正打开对应id页面
			if (fid == vfid) {
				//直接显示
				queryHistory(obj.from);
			} else {
				//修改未读消息小红点
				var count = parseInt($('#unread_' + vfid.toString()).text());
				$('#unread_' + vfid.toString()).text((++count).toString());
				$('#unread_' + vfid.toString()).show();
			}
			break;
		case MSG_VIDEO:
		
			var video_type = value.type;

			if (video_type == 'offer') {
				//bug
				fid = obj.from;
				
				navigator.getUserMedia({
					audio: true, // 是否开启麦克风
					video: true // 是否开启摄像头，这里还可以进行更多的配置
				}, function (stream) {
					// 绑定本地媒体流到video标签
					document.getElementById('localVideo').src = URL.createObjectURL(stream);
					// 向PeerConnection中加入需要发送的流
					remotePC.addStream(stream);
					remoteStream = stream;
					console.log('remote add stream');

					//保证getUserMedia之后才createAnswer
					var desc = new RTCSessionDescription(value);
					remotePC.setRemoteDescription(desc);
					remotePC.createAnswer(function (answer) {
						remotePC.setLocalDescription(answer);
						var msg = {
							"type": MSG_VIDEO,
							"from": parseInt($.cookie('uid')),
							"to": fid,
							"value": answer
						};
						ws.send(JSON.stringify(msg));
					}, errorHandler);
				}, function (error) {
					// 获取本地视频流失败
					alert("获取本地视频流失败", error);
				});

				//getUserMedia = (navigator.getUserMedia || navigator.webkitGetUserMedia || navigator.mozGetUserMedia || navigator.msGetUserMedia);
			} else if (video_type == 'answer') {
				//pc.setRemoteDescription(value);
				console.log('set remotedesc');
				pc.setRemoteDescription(new RTCSessionDescription(value));
			} else if (video_type == 'candidate_offer') {
				console.log('candidate_offer');
				console.log(value.value);
				remotePC.addIceCandidate(new RTCIceCandidate(value.value));
				/*
				remotePC.addIceCandidate(new RTCIceCandidate(value.value), function (err) {
					console.log('addice success' + err);
				}, function (err) {
					console.log('error:' + err);
				});
				*/
			} else if (video_type == 'candidate_anwser') {
				console.log('candidate_anwser');
				console.log(value.value);
				pc.addIceCandidate(new RTCIceCandidate(value.value));
				/*
				pc.addIceCandidate(new RTCIceCandidate(value.value), function (err) {
					console.log('addice success' + err);
				}, function (err) {
					console.log('error' + err);
				});
				*/
			}
			$('#video-chat').modal('show');
			break;
		case MSG_REJECT:
			//插入历史表与别的不一样
			addMsg(obj.type, obj.from, obj.to, value, MSG_RCV);
			break;
		default:
			alert("收到不支持的消息类型" + type);
		}
	}

	//获取未读消息
	$.post(
		'/unread',
		'uid=' + $.cookie('uid'),
		function (data) {
		var obj = JSON.parse(data);
		if (obj.num > 0) {
			for (var i = 0; i < obj.num; i++) {
				//插入历史表
				addMsg(obj.value[i].type, obj.value[i].to, obj.value[i].from, obj.value[i].value, MSG_RCV);

				//修改未读消息小红点
				var old = parseInt($('#unread_' + obj.value[i].from.toString()).text());
				$('#unread_' + obj.value[i].from.toString()).text((++old).toString());
				$('#unread_' + obj.value[i].from.toString()).show();
			}
		} else if (obj.num < 0) {
			alert("获取未读消息失败");
		}
	});
}

//发送消息
function sendMsg() {
	//获取消息内容
	txt_msg = $.trim(($('#right-edit').val()));
	//显示已发送的内容
	var p = document.createElement('p');
	p.innerHTML = "我：" + txt_msg;
	document.getElementById('right-show').appendChild(p);
	//发送消息
	var str = {
		"type": MSG_TEXT,
		"from": parseInt($.cookie('uid')),
		"to": fid,
		"value": txt_msg
	}
	ws.send(JSON.stringify(str));
	//记录历史表
	addMsg(MSG_TEXT, $.cookie('uid'), fid, $.trim(txt_msg), MSG_SEND);
	
	//bug
	//$('#right-edit').val('');
	$('.emoji-wysiwyg-editor').text('');
}
