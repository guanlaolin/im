var fid = 0; //记录好友id
var ws; //websocket
var pc; //RTCPeerConnection
var remotePC;
var db; //websql
var localStream; //控制摄像头关闭
var remoteStream;
//var getUserMedia;

//判断发送还是接收
var MSG_SEND = 0;
var MSG_RCV = 1;

//消息类型
var MSG_SUCCESS = 0;
var MSG_ERROR = 1;
var MSG_TEXT = 2;
var MSG_FILE = 3;
var MSG_VIDEO = 4;
var MSG_IDENTITY = 5;
var MSG_REJECT = 6;

//传输音视频
var options = {
	offerToReceiveAudio: true,
	offerToReceiveVideo: true
}

//绑定按键
onkeydown = function (evt) {
	if (evt.ctrlKey && evt.keyCode == 13) {
		//判断是否打开了聊天界面
		if ($('#foot-right').css('display') == 'block') {
			sendMsg();
		};
	}
}

$(function () {
	$('#video-chat').on('hide.bs.modal', function () {
		if (localStream){
			localStream.getVideoTracks()[0].stop();
			$('#localVideo').attr("src", "");
		}
		if (remoteStream) {
			remoteStream.getVideoTracks()[0].stop();
			$('#remoteVideo').attr("src", "");
		}
	});
})

//判断当前聊天好友
function setFid(obj) {
	//显示右侧
	$('#foot-right').show();
	//获取好友昵称
	var name = $('#' + obj).contents().filter(function () {
			return this.nodeType === 3;
		}).text();
	$('#right-title-name').text($.trim(name));
	//隐藏未读消息小红圈，并置0
	$('#unread_' + obj).hide();
	$('#unread_' + obj).text('0');

	fid = obj;

	queryHistory(obj);
}

//关闭按钮
function closeRight() {
	$('#foot-right').hide();
}

//登录
function login() {
	var uid = $.trim($('#txt_uid').val());
	var password = $('#txt_password').val();
	var remember = $('#cb_remember')[0].checked;

	//uid校验
	var checkResult = checkUid(uid);
	if (checkResult != true) {
		$('#login-error').text(checkResult);
		$('#login-error').show();
		return;
	}

	//密码校验
	checkResult = checkPassword(password);
	if (checkResult != true) {
		$('#login-error').text(checkResult);
		$('#login-error').show();
		return;
	}

	var value = 'uid=' + uid + '&password=' + password + '&remember=' + remember;

	$.post(
		'/login',
		value,
		function (data) {
		if (data == "success") {
			window.location.href = "/";
		} else {
			$('#login-error').text(data);
			$('#login-error').show();
		}
	});
}

//查找用户
function searchUser() {
	var uid = $.trim($('#search-uid').val());

	var checkResult = checkUid(uid);
	if (checkResult != true) {
		$('#search-alert').text(checkResult);
		$('#search-alert').show();
		return;
	}
	$.post(
		'/search',
		'fid=' + uid,
		function (data) {
		if ('empty' == data) {
			$('#search-alert').text('未找到用户');
			$('#tb-friend-list').hide();
			$('#search-alert').show();
		} else {
			var obj = JSON.parse(data);
			$('#user-list-id').text(obj.uid);
			$('#user-list-name').text(obj.uname);
			$('#user-list-func').attr('onclick', 'addFriend(' + obj.uid + ')');
			$('#search-alert').hide();
			$('#tb-friend-list').show();
		}
	});
}

//查看用户信息
function userInfo(uid) {
	$.post(
		'/userinfo',
		'fid=' + uid,
		function (data) {
			if ('empty' != data){
			var obj = JSON.parse(data);
				$('#info-uid').val(obj.uid);
				$('#info-uname').val(obj.uname);
				$('#info-email').val(obj.email);
				$('#info-motto').text(obj.motto);
			}
	});
	$('#user-info').modal('show');
}

//添加好友
function addFriend(fid) {
	$.post(
		'/addfriend',
		'fid=' + fid,
		function (data) {
		if (data != 'success') {
			$('#search-alert').text(data);
			$('#search-alert').show();
		} else {
			$('#search-alert').text('添加好友成功');
			$('#search-alert').show();

			window.location.href = '/';
			/*
			//更新好友列表
			var li = document.createElement('li')
			li.setAttribute('id',fid);
			li.setAttribute('onclick','setFid('+fid+')');
			li.innerHTML = $('#user-list-name').text();
			var span = document.createElement('span');
			span.setAttribute('class','unread_'+fid);
			span.innerHTML = 0;
			li.appendChild(span)
			document.getElementById('friend-list').appendChild(li);
			 */
		}
	});
}

//删除好友
function deleteFriend(fid) {
	$.post(
		'/delete',
		'fid=' + fid,
		function (data) {
		if ('success' == data) {
			$('#' + fid).remove();
			alert('成功删除好友' + fid);
		} else {
			alert(data);
		}
	});
	//清空好友对应的聊天记录
	deleleHistoryFromFid(fid);
}

//修改密码
//数据校验不全
function updatePsw(method) {
	var new_password = $('#new-password').val();
	var new_re_password = $('#new-re-password').val();
	var uid;
	var str;

	var checkResult = checkPassword(new_password);
	if (checkResult != true) {
		$('#update-psw-alert').text(checkResult);
		$('#update-psw-alert').show();
		return;
	}

	switch (method) {
	case 'reset': //重置密码
		uid = $.trim($('#uid').val());
		checkResult = checkUid(uid);
		if (checkResult != true) {
			$('#update-psw-alert').text(checkResult);
			$('#update-psw-alert').show();
			return;
		}
		var email = $.trim($('#email').val());
		checkResult = checkEmail(email);
		if (checkResult != true) {
			$('#update-psw-alert').text(checkResult);
			$('#update-psw-alert').show();
			return;
		}
		str = 'uid=' + uid + '&newpassword=' + new_password + '&newrepassword=' + new_re_password + '&email=' + email + '&type=' + method;
		break;
	case 'change': //修改密码
		uid = $.cookie('uid');
		var old_password = $('#old-password').val();
		checkResult = checkPassword(old_password);
		if (checkResult != true) {
			$('#update-psw-alert').text(checkResult);
			$('#update-psw-alert').show();
			return;
		}
		str = 'uid=' + uid + '&newpassword=' + new_password + '&newrepassword=' + new_re_password + '&oldpassword=' + old_password + '&type=' + method;
		break;
	default:
		alert('修改密码失败');
	}
	$.post(
		'/updatepsw',
		str,
		function (data) {
		if (data != 'success') {
			$('#update-psw-alert').text(data);
			$('#update-psw-alert').show();
		} else {
			$('#update-psw-alert').text('修改密码成功，请牢记您的新密码');
			$('#update-psw-alert').show();
		}
	});
}

//退出登录
function logout() {
	$.get(
		'/logout',
		'uid=' + $.cookie('uid'),
		function (data) {
		if (data == "success") {
			//清除cookie
			$.cookie('uid', '', {
				expires: -1
			});
			//清除ws
			ws = null;
			//跳转链接
			window.location.href = '/login';
		} else {
			alert(data);
		}
	});
}

//增加信息到历史表
function addMsg(type, from, to, value, flag) {
	db.transaction(function (ex) {
		ex.executeSql("INSERT INTO tb_history VALUES(?,?,?,?,?);", [type, from, to, value, flag]);
	});
}

//删除fid对应的历史记录
function deleleHistoryFromFid(fid) {
	db.transaction(function (ex) {
		ex.executeSql('DELETE FROM tb_history WHERE "to" = ?', [fid]);
	});
	$('#right').hide();
}

//历史记录
function queryHistory(fid) {
	$('#right-show').empty();
	db.transaction(function (ex) {
		ex.executeSql('SELECT * FROM tb_history WHERE "to" = ?', [fid],
			function (ex, results) {
			for (i = 0; i < results.rows.length; i++) {
				var p = document.createElement('p');
				var type = results.rows.item(i).type;
				var from = results.rows.item(i).from;
				var to = results.rows.item(i).to;
				var value = results.rows.item(i).value;
				var flag = results.rows.item(i).flag;
				if (flag == MSG_RCV) {
					p.innerHTML = $('#' + to).contents().filter(function () {
							return this.nodeType === 3;
						}).text() + ':' + value;
				} else {
					p.innerHTML = "我:" + value;
				}
				document.getElementById('right-show').appendChild(p);
			}
		}, null);
	});
}

function clearHistory(){
	db.transaction(function (ex) {
		ex.executeSql('DELETE FROM tb_history WHERE "to" = ?', [fid],null, null);
	});
	$('#right-show').empty();
}

//校验用户名格式
function checkUid(uid) {
	if (uid.length == 0 || uid.length > 11)
		return "uid长度必须为1-11";
	return true;
}

//校验密码格式
function checkPassword(password) {
	if (password.length == 0)
		return "密码不能为空";
	return true;
}

//校验邮箱格式
function checkEmail(email) {
	if (email.length == 0) {
		return "邮箱不能为空";
	}
	return true;
}
