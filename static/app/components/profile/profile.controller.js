'use strict';

var app = angular.module('ipl');

/**
 * Profile Controller
 * 
 * Controller for the profile page.
 */
app.controller('profileController', function($http, $window, urlService) {
    var vm = this;

    vm.init = init;
    // vm.userData = {
    //     firstName: 'Gal',
    //     lastName: 'Gadot',
    //     iNumber: 'I333333',
    //     alias: 'wonderwoman',
    //     points: 50,
    //     coins: 5,
    //     profilePic: '/static/assets/img/users/batman.jpeg'
    // };

    function init() {
        // console.log($http.defaults.headers.common.Authorization,urlService.userProfile)
        // $http.get(urlService.userProfile)
        var req = {
            url: '/api/profile',
            // headers: {
            //     Authorization: $window.localStorage.getItem('token')
            // },
            method: 'GET'
        }
        $http(req)
            .then(function(res) {
                // console.log('xxxxxxx',res.data, res.body);
                vm.userData = {
                    firstName: res.data.firstname,
                    lastName: res.data.lastname,
                    iNumber: res.data.inumber,
                    alias: res.data.alias,
                    // points: res.data.points,
                    coins: res.data.coin,
                    // profilePic: res.data.pic_loc
                };
                // $window.localStorage.displayName = 
                console.log('success');
            }, function(err) {
                console.log('error', err);
            });
    }
});
