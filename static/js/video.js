//视频聊天
function videoChat() {
	navigator.getUserMedia({
		audio: true, // 是否开启麦克风
		video: true // 是否开启摄像头，这里还可以进行更多的配置
	}, function (stream) {
		// 绑定本地媒体流到video标签
		document.getElementById('localVideo').src = URL.createObjectURL(stream);
		// 向PeerConnection中加入需要发送的流
		pc.addStream(stream);
		localStream = stream;
		console.log('pc addstream');

		//保证getUserMedia成功后才createOffer
		pc.createOffer(function (offer) {
			pc.setLocalDescription(offer, function () {
				console.log('setLocalDescription');
				var msg = {
					"type": MSG_VIDEO,
					"from": parseInt($.cookie('uid')),
					"to": fid,
					"value": offer
				};
				ws.send(JSON.stringify(msg));
				console.log(JSON.stringify(msg));
				console.log('send offer');
			}, errorHandler);
		}, errorHandler, options);
	}, function (error) {
		// 获取本地视频流失败
		alert("获取本地视频流失败" + error);
	});
}

//错误处理函数
function errorHandler(err) {
	alert(err);
}

function successHandler(err) {
	alert(err);
}
