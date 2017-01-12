(function(root, factory) {
    if ("function" === typeof define && define.amd) {
        define("tyle_system", [], factory);
    } else if ("object" === typeof module && module.exports) {
        module.exports = factory();
    } else {
        root.TileSystem = factory();
    }
})(this, function() {
    var TILE_SIZE = 256;

    function project(lat, lng) {
        var siny = Math.sin(lat * Math.PI / 180);
        siny = Math.min(Math.max(siny, -0.9999), 0.9999);
        return {
            x: TILE_SIZE * (0.5 + lng / 360),
            y: TILE_SIZE * (0.5 - Math.log((1 + siny) / (1 - siny)) / (4 * Math.PI))
        };
    }

    return {
        getTileXY: function(lat, lng, zoom) {
            var scale = 1 << zoom,
                worldCoords = project(lat, lng);
            return {
                x: Math.floor(worldCoords.x * scale / TILE_SIZE),
                y: Math.floor(worldCoords.y * scale / TILE_SIZE)
            };
        }
    };
});
