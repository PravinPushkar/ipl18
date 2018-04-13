'use strict';

var app = angular.module('ipl');

/**
 * Profile Controller
 * 
 * Controller for the profile page.
 */
app.controller('profileController', ['$http', '$window', '$rootScope', 'urlService', 'utilsService', function ($http, $window, $rootScope, urlService, utilsService) {
    var vm = this;

    var token;

    vm.init = init;

    // Init function for profile
    function init() {
        var currentUserINumber = $window.localStorage.getItem('iNumber');
        token = $window.localStorage.getItem('token');
        var params = {
            url: `${urlService.userProfile}/${currentUserINumber}`,
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Authorization': token
            }
        };
        $http(params)
            .then(function (res) {
                vm.userData = {
                    firstName: utilsService.capitalizeFirstLetter(res.data.firstname),
                    lastName: utilsService.capitalizeFirstLetter(res.data.lastname),
                    iNumber: utilsService.capitalizeFirstLetter(res.data.inumber),
                    alias: res.data.alias,
                    points: res.data.points,
                    coins: res.data.coin,
                    profilePic: res.data.picLocation
                };
                $window.localStorage.setItem('profilePic', vm.userData.profilePic);
                $window.localStorage.setItem('Name', vm.userData.firstName + " " +vm.userData.lastName);
            }, function (err) {
                if (err.data.code === 403 && err.data.message === 'token not valid') {
                    utilsService.logout('Session expired, please re-login', true);
                    return;
                }
                utilsService.showToast({
                    text: 'Error fetching profile',
                    hideDelay: 0,
                    isError: true
                });
            });
    }
}]);