var app = angular.module('ipl');

app.factory('socket', ["urlService", "$rootScope", function (urlService, $rootScope) {
  var webSocketUrl =
    ((window.location.protocol === "https:") ? "wss://" : "ws://") + window.location.host + urlService.feeds;
  var socket = new WebSocket(webSocketUrl);
  return {
    onopen: function (callback) {
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
    send: function (message) {
      socket.send(message);
    }
  };
}]);