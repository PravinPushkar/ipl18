'use strict';

var app = angular.module('ipl');

app.controller('teamsController', function ($http, urlService) {
    var vm = this;

    vm.init = init;

    vm.teamsList = [{
        name: 'csk'
    },{
        name: 'dd'
    }];

    function init() {
        var params = {
            url: urlService.teams,
            method: 'GET',
            headers: {
                'Accept': 'application/json'
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