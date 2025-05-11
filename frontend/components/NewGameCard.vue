<script setup>

const props = defineProps({
    game: {
        default: SAMPLE_GAME
    }
})

function newGameCardHover() {

}
</script>

<template>
    <div class="new-game-card" @mouseover="newGameCardHover">
        <img class="new-game-thumbnail" :src="game.thumbnail">
        <p class="new-game-title"> {{ game.title }} </p>
        <p class="new-game-description"> {{ game.description }} </p>
    </div>
</template>

<style lang="scss" scoped>
@use '~/assets/scss/colors';
@use '~/assets/scss/util';

$thumbnail-width: 50px;

.new-game-card {
    display: grid;
    grid-template-columns: $thumbnail-width auto;
    grid-template-rows: 1fr 1fr;
    gap: 3px;
    background-color: colors.$light-bg-color;
    border: 3px solid;
    border-color: navy;
}

.new-game-card:hover {
    border-color: lightseagreen;
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

.new-game-card a {
    text-decoration: none;
    color: black;
    border-style: solid;
}

.new-game-thumbnail {
    max-height: $thumbnail-width;
    max-width: $thumbnail-width;
}
</style>