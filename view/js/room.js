var isHomeowner = false;
var isPrepare = false


function joinRoom() {
    var roomId = $("#RoomId").val();
    var json = {
        Room: {
            UserName: User.userName,
            UserId: User.userId,
            RoomId: roomId,
        }
    }
    SendJson(json)
}

function SendMessage() {
    var SendMessage = $("#SendMessage").val();
    $("#SendMessage").val("")
    var json = {
        GameMessage: {
            UserName: User.userName,
            UserId: User.userId,
            RoomId: User.roomId,
            Msg: SendMessage
        }
    }
    SendJson(json)
}

function createRoom() {
    var totalNumber = $("#TotalNumber").val();
    var UndercoverNumber = $("#UndercoverNumber").val();
    var json = {
        Room:{
            UserName: User.userName,
            UserId: User.userId,
            RoomPassword: "",
            TotalNumber: totalNumber,
            UndercoverNumber: UndercoverNumber,
        }
    }
    isHomeowner = true
    SendJson(json)
}

function rRoomCallBack(data) {
    console.log("进入房间")
    roomInfo = data.GameMessage.Data.RoomInfo
    // roomId = roomInfo.RoomId
    // $.cookie('roomId' + code, roomId);
    // if (roomInfo.CreateUserId == User.userId) {
    //     isHomeowner = true
    // }

    initRoom(roomInfo.RoomId, roomInfo.CreateUserId)
}


function rPrepareCallBack(data) {
    console.log("准备" + isPrepare)
    console.log("准备")
    isPrepare = !isPrepare
    var btnHtml = "准备"
    if (isPrepare) {
        btnHtml = "准备中"
    }
    $(".btn-zb").html(btnHtml)
}

function rStartCallBack(data) {
    console.log("开始回调")
    if (data.GameMessage.Msg === "开始") {
        $(".btn-ks").removeClass("weui-btn_disabled")
    } else {
        $(".btn-ks").addClass("weui-btn_disabled")
    }
}

function Prepare() {
    console.log("准备")
    if (isHomeowner) {
        return
    }
    var roomId =  $.cookie('roomId' + code);
    var json = {
        Room:{
            UserName: User.userName,
            UserId: User.userId,
            RoomId: roomId,
            RoomPassword: "",
            IsPrepare: true,
        }
    }
    SendJson(json)
}

function StartGame() {
    console.log("开始");
    var b = $(".btn-ks").hasClass("weui-btn_disabled");
    if (b) {
        return
    }

    var json = {
        Game:{
            Stage: "Start",
            RoomId: User.roomId,
        }
    }
    isHomeowner = true
    SendJson(json)
}

function initRoom(roomId, CreateUserId) {
    if (roomId == "") {
        roomId = $.cookie('roomId' + code, roomId)
    } else {
        $.cookie('roomId' + code, roomId, { expires: 1 });
    }

    User.roomId = roomId

    if (CreateUserId == User.userId) {
        isHomeowner = true
    }

    isPrepare = false
    $("#page1").load("./room.html");
    $(".return").addClass("out");
}