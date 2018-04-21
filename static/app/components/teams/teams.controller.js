'use strict';

var app = angular.module('ipl');

app.controller('teamsController', ['$http', '$window', 'urlService', 'utilsService', function ($http, $window, urlService, utilsService) {
    var vm = this;

    var token;

    vm.init = init;

    function init() {
        token = $window.localStorage.getItem('token');
        var params = {
            url: urlService.teams,
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Authorization': token
            }
        };
        vm.teamsList = [];
        $http(params)
            .then(function (res) {
                vm.teamsList = res.data.teams;
            }, function (err) {
                if (err.data.code === 403 && err.data.message === 'token not valid') {
                    utilsService.logout('Session expired, please re-login', true);
                    return;
                }
                utilsService.showToast({
                    text: 'Error in fetching teams',
                    hideDelay: 0,
                    isError: true
                });
            });
    }
}]);