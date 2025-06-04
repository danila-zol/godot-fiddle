<script setup>
const props = defineProps({
    game: {
        default: SAMPLE_GAME2
    },
    thumbnailWidth: { default: "50px" },
    thumbnailHeight: { default: "50px" },
    gameId: { default: 1 }
})
</script>

<template>
    <NuxtLink :to="'/games/' + gameId" class="new-game-card">
        <img class="new-game-thumbnail" alt="Картинка игры" :width="thumbnailWidth" :height="thumbnailHeight"
            :src="game.thumbnail" fetchpriority=high>
        <p class="new-game-title"> {{ game.title }} </p>
        <p class="new-game-description"> {{ game.description }} </p>
    </NuxtLink>
</template>

<style lang="scss" scoped>
@use '~/assets/scss/colors';
@use '~/assets/scss/util';

.new-game-card {
    display: grid;
    grid-template-columns: auto 1fr;
    grid-template-rows: 24px 1fr;
    gap: 3px;
    background-color: colors.$base-color;
    border: 3px solid;
    border-color: colors.$light-highlight-color-darkererer;
    img {
        grid-row: span 2 ;
    }
}

.new-game-card:hover {
    border-color: colors.$light-highlight-color-darkerer;
}

.new-game-title {
    font-size: 18px;
    font-weight: bold;
    margin: 0px;
    padding: 0px;
    grid-row: 50px;
    align-self: center;
}

.new-game-description {
    font-size: 14px;
    margin: 0px;
    padding: 0px;
}

.new-game-card {
    text-decoration: none;
    color: black;
    border-style: solid;
}
</style>