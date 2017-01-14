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

    var TileSystem =  {
        getTileXY: function(lat, lng, zoom) {
            var scale = 1 << zoom,
                worldCoords = project(lat, lng);
            return {
                x: Math.floor(worldCoords.x * scale / TILE_SIZE),
                y: Math.floor(worldCoords.y * scale / TILE_SIZE)
            };
        },

        tileXYToQuadKey: function(x, y, zoom) {
            var buffer = [], digit, mask;
            for (var i = zoom; i > 0; --i) {
                digit = 0;
                mask = 1 << (i - 1);
                if ((x & mask) !== 0) {
                    ++digit
                }

                if ((y & mask) !== 0) {
                    digit += 2;
                }

                buffer.push(digit);
            }

            return buffer.join('');
        },

        getQuadKey: function(lat, lng, zoom) {
            var tile = TileSystem.getTileXY(lat, lng, zoom);
            return TileSystem.tileXYToQuadKey(tile.x, tile.y, zoom);
        }
    };

    return TileSystem;
});
