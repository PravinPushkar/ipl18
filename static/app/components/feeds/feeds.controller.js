'use strict';

var app = angular.module('ipl');

/**
 * feedsController Controller
 * 
 * Controller for feeds page
 */
app.controller('feedsController', ['$http', '$window', 'socket', '$timeout', '$scope', 'utilsService', '$rootScope', '$location', function ($http, $window, socket, $timeout, $scope, utilsService, $rootScope, $location) {
    var feeds = this;

    var token;
    feeds.feedEntries = [];
    feeds.init = init();
    feeds.submit = submit;
    var currentUserINumber = $window.localStorage.getItem('iNumber');
    feeds.buzz = "";
    var reg = /^[a-zA-Z0-9?!_\-, .*():]+$/;
    feeds.isLoaded = false;
    function init() {
        token = $window.localStorage.getItem('token');
        var auth = {
            'Authorization': token
        }
        socket.onopen(function () {
            socket.send(JSON.stringify(auth));
        });
        socket.onmessage(function (data) {
            $timeout(function () {
                feeds.isLoaded = true;
                $scope.$apply(function () {
                    feeds.feedEntries.push(JSON.parse(data));
                });
            });
        });
    }
    $rootScope.$on('$locationChangeStart', function () {
        if ($location.path() !== '/feeds') {
            socket.onclose();
        }
    });

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