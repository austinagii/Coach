<div class="loading-indicator"></div>

<style>
    .loading-indicator {
        animation: pulse 1s infinite; 
        animation-timing-function: ease-in-out;
        border: 1px solid black;
        border-radius: 100%;
    }
    
    @keyframes pulse {
        0% { color: black; height: 0em; width: 0em; }
        50% { border-color: white; }
        100% { color: black; height: 5em; width: 5em; }
    }
</style>