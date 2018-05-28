let app = require('lotion')({
  initialState: {  
  	'test': {
  		rank: 1,
  		links: {
  			'QmVWVscvorUwah9x9GL2465pLFT97GQxXtJaVZQXvsBvaJ': 1
  		}
  	}
  }
})

app.use((state, tx) => {
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
	// console.log('>> ', tx.type == 'link');
  // state.count++
})

app.listen(5000)