(function(root, factory) {
    if ("function" === typeof define && define.amd) {
        define("app", ["google_maps", "remote_manager", "tile_system"], factory);
    } else if ("object" === typeof module && module.exports) {
        module.exports = factory(
            require('google_maps'),
            require('remote_manager'),
            require('tile_system')
        );
    } else {
        root.App = factory(
            root.google,
            root.RemoteManager,
            root.TileSystem
        );
    }
})(this, function(google, RemoteManager, TileSystem) {

    var IMAGE_PATH = 'https://googlemaps.github.io/js-marker-clusterer/images/m',
        IMAGE_EXT = 'png',
        ICONS_COUNT = 5;

    function init(cfg) {
        var map = new google.maps.Map(document.getElementById(cfg.element), {
            center: {lat: 43.124365, lng: 131.889356},
            zoom: cfg.zoom
        });

        function CoordMapType() {
            this.tileSize = new google.maps.Size(256, 256);
        }
        
        CoordMapType.prototype.getTile = function(coord, zoom, ownerDocument) {
            var div = ownerDocument.createElement('div');
            div.innerHTML = coord + '; ' + TileSystem.tileXYToQuadKey(coord.x, coord.y, zoom);
            div.style.width = this.tileSize.width + 'px';
            div.style.height = this.tileSize.height + 'px';
            div.style.fontSize = '10';
            div.style.borderStyle = 'solid';
            div.style.borderWidth = '1px';
            div.style.borderColor = '#AAAAAA';
            return div;
        };

        map.overlayMapTypes.insertAt(0, new CoordMapType());

        var onPointsRecieved = function(data) {
            var sum = 0, cluster;

            for (var k in data) {
                cluster = new Cluster(data[k]);
                addMarker(cluster);
                sum += cluster.getCount();
            }

            console.log('ADDED '+ sum + ' CNT');
            console.log('TOTAL '+ markers.length + ' CNT');
        };

        var remoteManager = new RemoteManager(
            map, cfg.wsPath, onPointsRecieved
        );

        remoteManager.run();

        var Cluster = function(data) {
            this.lat = data.lat;
            this.lng = data.lng;
            this.count = data.count;
        };

        Cluster.prototype.getCount = function() {
            return this.count;
        };

        Cluster.prototype.getPosition = function() {
            return {lat: this.lat, lng: this.lng};
        };

        Cluster.prototype.getIcon = function() {
            var dv = this.count || 0,
            index = 0;

            while (dv !== 0 && index <= ICONS_COUNT) {
                dv = parseInt(dv / 10);
                index++;
            }

            return IMAGE_PATH + index + '.' + IMAGE_EXT;
        };

        Cluster.prototype.getLabel = function() {
            return this.count.toString();
        };

        Cluster.prototype.asMarker = function() {
            return new google.maps.Marker({
                position: this.getPosition(),
                map: map,
                label: this.getLabel(),
                icon: this.getIcon()
            });
        }

        var markers = [];

        var addMarker = function(cluster) {
            markers.push(cluster.asMarker());
        };

        function setMapOnAll(map) {
            for (var i = 0; i < markers.length; i++) {
                markers[i].setMap(map);
            }
        }

        var clearAllMarkers = function() {
            setMapOnAll(null);
        };

        var deleteAllMarkers = function() {
            clearAllMarkers(null);
            markers = [];
        };

        map.addListener('zoom_changed', function() {
            deleteAllMarkers();
        });
    }

    return {
        init: init
    };
});
