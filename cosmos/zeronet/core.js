

const core = (state, tx) => {
    const { type, keyword, hash } = tx;

    if (type === 'search') {
        if (state[keyword]) {
            state[keyword].rank++;
        } else {
            state[keyword] = {
                rank: 1,
                links: {}
            }
        }
    } 

    if (type === 'link') {
        if (state[keyword]) {
            if (state[keyword].links[hash]) {
                state[keyword].links[hash]++
            } else {
                state[keyword].links[hash] = 1
            }
        } else {
            state[keyword] = {
                rank: 1,
                links: {}
            }
        }
    }
}

module.exports = core;