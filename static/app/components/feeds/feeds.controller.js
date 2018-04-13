'use strict';

var app = angular.module('ipl');

/**
 * feedsController Controller
 * 
 * Controller for feeds page
 */
app.controller('feedsController', ['$http', '$window', 'socket', '$timeout', '$scope', 'utilsService', function ($http, $window, socket, $timeout, $scope, utilsService) {
    var feeds = this;

    var token;

    feeds.init = init();
    feeds.feedEntries = [];
    feeds.submit = submit;
    var currentUserINumber = $window.localStorage.getItem('iNumber');
    feeds.buzz = "";
    var reg = /^[a-zA-Z0-9?!_\-, .*()]+$/;
    function init() {
        token = $window.localStorage.getItem('token');
        var auth = {
            'Authorization': token
        }
        socket.onopen(function () {
            socket.send(JSON.stringify(auth));
            console.log("connection established");
        });
        socket.onmessage(function (data) {
            $timeout(function () {
                $scope.$apply(function () {
                    feeds.feedEntries.push(JSON.parse(data));
                });
            });
        });
    }

    function submit() {
        if (feeds.buzz && reg.test(feeds.buzz.trim())) {
            socket.send(feeds.buzz)
            feeds.buzz = "";
        } else {
            utilsService.showToast({
                text: 'Type something worthy!',
                hideDelay: 2000,
                isError: true
            });
        }
    }
    socket.onerror(function (err) {
        utilsService.showToast({
            text: 'Error in websocket connect, Please refresh the page',
            hideDelay: 2000,
            isError: true
        });
    });

}]);