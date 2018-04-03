'use strict';

var app = angular.module('ipl');

/**
 * Edit Profile Controller
 * 
 * Controller for the edit profile page.
 */
app.controller('editProfileController', ['$http', '$mdToast', '$scope', '$state', '$window', 'urlService', 'utilsService', 'aliasPattern', function ($http, $mdToast, $scope, $state, $window, urlService, utilsService, aliasPattern) {
    var vm = this;

    var token;

    vm.edit = edit;
    vm.init = init;

    vm.aliasPattern = aliasPattern;

    // Get the file name of image uploaded
    document.querySelector('input[id="profilePic"]').onchange = function () {
        vm.imageSelectedName = document.getElementById('profilePic').files[0].name;
    };

    // $scope.$on('fileProgress', function (e, progress) {
    //     vm.progress = progress.loaded / progress.total;
    // });

    // Init function for edit profile
    function init() {
        token = $window.localStorage.getItem('token');
    }

    // Send editable data to back-end
    function edit() {
        if ((vm.password !== '' && vm.password !== undefined && vm.password !== null) && (vm.alias !== '' && vm.alias !== undefined && vm.alias !== null) && !(document.getElementById('profilePic').files[0])) {
            utilsService.showToast({
                text: 'Please enter valid value in the fields.',
                hideDelay: 2000,
                isError: true
            });
            return;
        }
        if (vm.password !== vm.confirmPassword) {
            utilsService.showToast({
                text: 'Password and Confirm Password do not match',
                hideDelay: 2000,
                isError: true
            });
            return;
        }
        var fd = new FormData();
        if (vm.password !== '' && vm.password !== undefined && vm.password !== null) {
            fd.append('password', vm.password);
        }
        if (vm.alias !== '' && vm.alias !== undefined && vm.alias !== null) {
            fd.append('alias', vm.alias);
        }
        if (document.getElementById('profilePic').files[0]) {
            fd.append('image', document.getElementById('profilePic').files[0]);
        }
        var currentUserINumber = $window.localStorage.getItem('iNumber');

        var params = {
            url: `${urlService.userProfile}/${currentUserINumber}`,
            data: fd,
            method: 'PUT',
            transformRequest: angular.identity,
            headers: {
                'Content-Type': undefined,
                'Authorization': token
            }
        };
        $http(params)
            .then(function () {
                utilsService.showToast({
                    text: 'User Profile updated.',
                    hideDelay: 1500,
                    isError: false
                });
                $state.go('main.profile');
                return;
            }, function (err) {
                if (err.data.code === 403 && err.data.message === 'token not valid') {
                    utilsService.logout('Session expired, please re-login', true);
                    return;
                }
                utilsService.showToast({
                    text: 'Error in updating profile',
                    hideDelay: 2000,
                    isError: true
                });
                return;
            });
    }
}]);