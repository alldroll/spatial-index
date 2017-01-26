(function(root, factory) {
    if ("function" === typeof define && define.amd) {
        define("inherence", [], factory);
    } else if ("object" === typeof module && module.exports) {
        module.exports = factory();
    } else {
        root.Inherence = factory();
    }
})(this, function() {
    return function(child, parent) {
        for (var k in parent) {
            if (parent.hasOwnProperty(k)) {
                child[k] = parent[k];
            }
        }

        function Constructor() {
            this.constructor = child;
        }

        Constructor.prototype = parent.prototype;
        child.prototype = new Constructor();
        child.__super__ = parent.prototype;

        return child;
    };
});
