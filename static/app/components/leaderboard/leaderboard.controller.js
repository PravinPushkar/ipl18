'use strict';

var app = angular.module('ipl');

/**
 * Leaderboard Controller
 * 
 * Controller for leaderboard page
 */
app.controller('leaderboardController', ['$http', '$window', 'utilsService', 'urlService', function ($http, $window, utilsService, urlService) {
    var vm = this;

    var token;

    vm.init = init();

    // vm.leaderboardData = [{
    //     firstName: utilsService.capitalizeFirstLetter('Akshil'),
    //     lastName: utilsService.capitalizeFirstLetter('verma'),
    //     alias: 'abalrok',
    //     iNumber: utilsService.capitalizeFirstLetter('i341668'),
    //     points: 55,
    //     profilePic: '/static/assets/img/users/batman.jpeg'
    // }, {
    //     firstName: utilsService.capitalizeFirstLetter('gal'),
    //     lastName: utilsService.capitalizeFirstLetter('gadot'),
    //     alias: 'wonderwoman',
    //     iNumber: utilsService.capitalizeFirstLetter('i313131'),
    //     // points: parseInt('22'),
    //     points: 22,
    //     profilePic: '/static/assets/img/users/galgadot.jpg'
    // }];
    // var points = [];
    // vm.leaderboardData.forEach(function (user) {
    //     points.push(parseInt(user.points));
    // });
    // vm.highestPoints = Math.max(...points);

    // Init function for the leaderboard view
    function init() {
        token = $window.localStorage.getItem('token');
        var params = {
            url: urlService.leaderboard,
            method: 'GET',

            headers: {
                'Accept': 'application/json',
                'Authorization': token
            }
        };
        $http(params)
            .then(function (res) {
                var points = [];
                vm.leaderboardData = [];
                res.data.leaders.forEach(function(user){
                vm.leaderboardData.push({
                    firstName: utilsService.capitalizeFirstLetter(user.firstname),
                    lastName: utilsService.capitalizeFirstLetter(user.lastname),
                    iNumber: utilsService.capitalizeFirstLetter(user.inumber),
                    alias: user.alias,
                    points: parseInt(user.points),
                    //coins: user.coin,
                    profilePic: user.piclocation
                });
                points.push(parseInt(user.points));
                });
                vm.highestPoints = Math.max(...points);
                console.log('success');
            }, function (err) {
                if (err.data.code === 403 && err.data.message === 'token not valid') {
                    utilsService.logout('Session expired, please re-login', true);
                    return;
                }
            });
    }

}]);