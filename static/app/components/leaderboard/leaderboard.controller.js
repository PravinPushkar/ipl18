'use strict';

var app = angular.module('ipl');

/**
 * Leaderboard Controller
 * 
 * Controller for leaderboard page
 */
app.controller('leaderboardController', function ($http, utilsService, urlService) {
    var vm = this;

    vm.init = init();

    vm.leaderboardData = [{
        firstName: utilsService.capitalizeFirstLetter('Akshil'),
        lastName: utilsService.capitalizeFirstLetter('verma'),
        alias: 'abalrok',
        iNumber: utilsService.capitalizeFirstLetter('i341668'),
        points: 55,
        profilePic: '/static/assets/img/users/batman.jpeg'
    }, {
        firstName: utilsService.capitalizeFirstLetter('gal'),
        lastName: utilsService.capitalizeFirstLetter('gadot'),
        alias: 'wonderwoman',
        iNumber: utilsService.capitalizeFirstLetter('i313131'),
        // points: parseInt('22'),
        points: 22,
        profilePic: '/static/assets/img/users/galgadot.jpg'
    }];
    var points = [];
    vm.leaderboardData.forEach(function(user) {
        points.push(parseInt(user.points));
    });
    vm.highestPoints = Math.max(...points);

    function init() {
        var params = {
            url: urlService.leaderboard,
            method: 'GET',
            headers: {
                'Accept': 'application/json'
            }
        };
        $http(params)
            .then(function (res) {
                // var points = [];
                // vm.leaderboardData = [];
                // res.data.forEach(function(user){
                // vm.leaderboardData.push({
                //     firstName: utilsService.capitalizeFirstLetter(user.firstname),
                //     lastName: utilsService.capitalizeFirstLetter(user.lastname),
                //     iNumber: utilsService.capitalizeFirstLetter(user.inumber),
                //     alias: user.alias,
                //     // points: parseInt(user.points),
                //     coins: user.coin,
                //     // profilePic: user.piclocation
                // });
                // points.push(parseInt(user.points));
            // });
            // vm.highestPoints = Math.max(...points);
                console.log('success');
            }, function () {
                console.log('error');
            });
    }

});