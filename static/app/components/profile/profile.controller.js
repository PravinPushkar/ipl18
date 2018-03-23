'use strict';

var app = angular.module('ipl');

/**
 * Profile Controller
 * 
 * Controller for the profile page.
 */
app.controller('profileController', function ($http, $window, urlService, utilsService) {
    var vm = this;

    vm.init = init;

    vm.setAlias = $window.localStorage.getItem('setAlias');

    vm.userData = {
        firstName: 'bruce',
        lastName: 'wayne',
        iNumber: 'I333333',
        alias: 'chamgadar_aaaamaanav',
        points: 50,
        coins: 5,
        profilePic: '/static/assets/img/users/batman.jpeg'
    };

    function init() {
        var params = {
            url: urlService.userProfile,
            method: 'GET',
            headers: {
                'Accept': 'application/json'
            }
        };
        $http(params)
            .then(function (res) {
                vm.userData = {
                    firstName: utilsService.capitalizeFirstLetter(res.data.firstname),
                    lastName: utilsService.capitalizeFirstLetter(res.data.lastname),
                    iNumber: utilsService.capitalizeFirstLetter(res.data.inumber),
                    alias: res.data.alias,
                    // points: res.data.points,
                    coins: res.data.coin,
                    // profilePic: res.data.piclocation
                };
                $window.localStorage.setItem('displayName', vm.setAlias ? vm.userData.alias : `${vm.userData.firstName} ${vm.userData.lastName}`);
                $window.localStorage.getItem('setAlias', vm.setAlias);
                console.log('success');
            }, function(err) {
                console.log('error', err);
            });
    }
});
