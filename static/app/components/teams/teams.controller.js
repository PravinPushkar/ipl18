'use strict';

var app = angular.module('ipl');

app.controller('teamsController', function ($http, $window, urlService) {
    var vm = this;

    var token;
    
    vm.init = init;

    vm.teamsList = [{
        name: 'csk'
    },{
        name: 'dd'
    }];

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
        vm.teams = [];
        $http(params)
            .then(function (res) {
                console.log('success');
                // res.data.forEach(function(team) {
                //     vm.teams.push({
                //         name: team.name,
                //         alias: team.shortname
                //     });
                // });
            }, function () {
                console.log('error');
            });
    }
});