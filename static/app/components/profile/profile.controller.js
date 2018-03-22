'use strict';

var app = angular.module('ipl');

/**
 * Profile Controller
 * 
 * Controller for the profile page.
 */
app.controller('profileController', function ($http, $window, urlService) {
    var vm = this;

    vm.init = init;
    vm.userData = {
        firstName: 'Gal',
        lastName: 'Gadot',
        iNumber: 'I333333',
        alias: 'wonderwoman',
        points: 50,
        coins: 5,
        profilePic: '/static/assets/img/users/batman.jpeg'
    };

    function init() {
        $http.get(urlService.userProfile)
            .then(function (res) {
                vm.userData = {
                    firstName: res.body.firstName,
                    lastName: res.body.lastName,
                    iNumber: res.body.iNumber,
                    alias: res.body.alias,
                    points: res.body.points,
                    coins: res.body.coins,
                    // profilePic: res.body.pic_loc
                };
                $window.localStorage.displayName = 
                console.log('success');
            }, function () {
                console.log('error');
            });
    }
});