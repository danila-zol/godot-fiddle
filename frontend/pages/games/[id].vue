<script setup>
const runConf = useRuntimeConfig()
const route = useRoute()
const gameId = route.params.id

let engineBinaryPath = '/thirdparty/godot/godot.editor'
let engineJsPath = '/thirdparty/godot/godot.editor.js'

let demoUrl = runConf.public.apiRoot + runConf.public.apiDemosPrefix + "/" + gameId
let { data: demo, error: error } = await useFetch(demoUrl)
let demoZip
let gameFsPathPrefix = '/tmp'
let gameFsFilename
let gameTitle
let gameDescription

if (true) {
    gameFsFilename = "dodge_the_creeps.zip"
    gameTitle = SAMPLE_GAME2.title
    gameDescription = SAMPLE_GAME2.description
    demoZip = await fetch("/thirdparty/godot/dodge_the_creeps.zip")
        .then((res) => res.arrayBuffer())
}

let gameFsPath = gameFsPathPrefix + "/" + gameFsFilename
let args = "--main-pack " + gameFsPath
let gameConfig = {
    'unloadAfterInit': false,
    'canvasResizePolicy': 1,
    'mainPack': "/thirdparty/godot/dodge_the_creeps.zip",
    'executable': engineBinaryPath
}

let game

let loadEngine = () => {
    game = new Engine(gameConfig)
    game.startGame()
}

useHead({
    script: [
        {
            hid: 'godot',
            src: engineJsPath,
            defer: true,
            async: true,
            onload: loadEngine
        }
    ]
})

definePageMeta({
    validate: async (route) => {
        return typeof route.params.id === 'string' && /^\d+$/.test(route.params.id)
    }
})
</script>

<template>
    <Navbar></Navbar>
    <div id="centering-containter">
       <p class="game-title" style="font-size: 34px;">{{ gameTitle }}</p>
        <div id="game-area">
            <canvas id="game-canvas" width="800" height="600"></canvas>
            <p style="font-size: 24px;font-weight: bold;width:95%;margin: 0 10px;">Описание</p>
            <p class="game-description">{{ gameDescription }}</p>
        </div>
    </div>
</template>

<style lang="scss" scoped>
.tab-hidden {
    display: none;
}

#game-title-input {
    border: none;
    border-bottom: solid 2px black;
    font-size: 28px;
    font-weight: bold;
    margin: 3px 10px;
}

#site-actions {
    margin-left: 50px;
}

#centering-containter {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
}

#game-area {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    max-width: 80%;
    background-color: grey;
}

#game-area canvas {
    margin: 10px 0
}

.game-description {
    max-width: 800px;
    margin: 0 10px;
    font-size: 18px;
}
</style>