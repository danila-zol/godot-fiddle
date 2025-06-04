<script setup>
const props = defineProps({
    game: {
        default: SAMPLE_GAME
    },
    thumbnailWidth: { default: "50px" },
    thumbnailHeight: { default: "50px" },
    gameId: { default: 1 }
})
let maxDescriptionLen = 170
let gameDescriptionCropped = computed(() => {
    let gameDesc = props.game.description
    if (gameDesc.length < maxDescriptionLen) {
        return gameDesc
    } else {
        return gameDesc.slice(0, maxDescriptionLen) + " ..."
    }
})
</script>

<template>
    <NuxtLink :to="'/games/' + gameId" class="new-game-card">
        <img class="new-game-thumbnail" alt="Картинка игры" :width="thumbnailWidth" :height="thumbnailHeight"
            :src="game.thumbnail" fetchpriority=high>
        <p class="new-game-title"> {{ game.title }} </p>
        <p class="new-game-description"> {{ gameDescriptionCropped }} </p>
    </NuxtLink>
</template>

<style lang="scss" scoped>
@use '~/assets/scss/colors';
@use '~/assets/scss/util';

$thumbnail-width: 50px;

.new-game-card {
    display: grid;
    grid-template-columns: auto 1fr;
    grid-template-rows: auto 1fr;
    gap: 3px;
    background-color: colors.$base-color;
    border: 3px solid;
    border-color: colors.$light-highlight-color-darkererer;
}

.new-game-card:hover {
    border-color: colors.$light-highlight-color-darkerer;
}

.new-game-title {
    font-size: 18px;
    font-weight: bold;
    margin: 0px;
    padding: 0px;
}

.new-game-description {
    font-size: 14px;
    margin: 0px;
    padding: 0px;
    grid-column: span 2;
}

.new-game-card {
    text-decoration: none;
    color: black;
    border-style: solid;
}
</style>