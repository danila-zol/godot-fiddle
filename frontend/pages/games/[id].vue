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

onMounted(() => {
    if (typeof Engine == 'function') {
        loadEngine()
    }
})

definePageMeta({
    validate: async (route) => {
        return typeof route.params.id === 'string' && /^\d+$/.test(route.params.id)
    }
})

// Cleanup
onBeforeRouteLeave(() => {
    if (game) {
        game.requestQuit();
        game.value = undefined
    }
})
</script>

<template>
    <Navbar></Navbar>
    <div class="centering-containter">
        <p class="game-title" style="font-size: 34px;">{{ gameTitle }}</p>
        <div id="game-area">
            <p v-if="!game.value">Godot загружается</p>
            <canvas id="game-canvas" width="800" height="600"></canvas>
            <p style="font-size: 24px;font-weight: bold;width:95%;margin: 0 10px;">Описание</p>
            <p class="game-description">{{ gameDescription }}</p>
        </div>
    </div>
</template>

<style lang="scss" scoped>
@use '~/assets/scss/_colors';

.tab-hidden {
    display: none;
}

#game-title-input {
    border: none;
    border-bottom: solid 2px colors.$light-highlight-color-darker;
    font-size: 28px;
    font-weight: bold;
    margin: 3px 10px;
}

#site-actions {
    margin-left: 50px;
}

#game-area {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    max-width: 80%;
    background-color: colors.$base-color;
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