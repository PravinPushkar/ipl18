var app = angular.module('ipl');

app.factory('socket', ["urlService", "$rootScope", function (urlService, $rootScope) {
  var webSocketUrl =
    ((window.location.protocol === "https:") ? "wss://" : "ws://") + window.location.host + urlService.feeds;
  var socket;
  return {
    onopen: function (callback) {
      socket = new WebSocket(webSocketUrl);
      socket.onopen = function () {
        $rootScope.$apply(function () {
          callback.apply(socket);
        });
      };
    },
    onmessage: function (callback) {
      socket.onmessage = function (e) {
        var data = e.data;
        $rootScope.$apply(function () {
          if (callback) {
            callback.apply(socket, [data]);
          }
        });
      }
    },
    onerror: function (callback) {
      socket.onerror = function (error) {
        $rootScope.$apply(function () {
          if (callback) {
            callback.apply(socket, [error]);
          }
        });
      };
    },
    onclose: function () {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close();
      }
    },
    send: function (message) {
      socket.send(message);
    }
  };
}]);