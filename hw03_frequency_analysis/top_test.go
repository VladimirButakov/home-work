package hw03frequencyanalysis

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = false

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

var textEnglish = `By 1999, the Cold War had thawed, and it seemed 
	nuclear proliferation would soon be a thing of the past. Despite 
	this, all was not well in the world. A series of shocks to the oil 
	market spurred the development of new high-tech energy sources,
	including fusion power. However, most vehicles still relied on oil. 
	Oil reserves were at a critical low, and the world community was 
	prepared to take drastic measures, either by drilling into sand and 
	shale for more oil, despite the difficulty—or moving on to renewable 
	fuels.
	Such steps proved unnecessary when Czech scientist, Dr. Kio Marv, 
	successfully bio-engineered a new species of algae, 
	OILIX, that could produce petroleum-grade hydrocarbons with little 
	expense and effort. Marv was on his way to a demonstration in the United 
	States when he was kidnapped by soldiers from Zanzibar Land. NATO 
	discovered that Zanzibar Land's leaders planned to hold the world hostage 
	by controlling the supply of oil, and some good old-fashioned nuclear 
	brinkmanship, courtesy of a stockpile of nukes.
	Solid Snake was brought out of retirement by FOXHOUND's new commander, 
	Roy Campbell, and was sent to Zanzibar Land to rescue Dr. Marv.
	For a full summary, see Zanzibar Land Disturbance.
	The intro to the game and the instruction manual mention that nuclear 
	weapons had been completely abandoned by the time of the main plot, 
	making Zanzibar Land the world's sole nuclear power. Metal Gear Solid, 
	however, retcons this account by having the current nuclear-armed nations 
	maintain their stockpiles, with the reduction of nuclear weapons via the 
	START-3 Treaty serving a prominent role in the story. Any references 
	to global nuclear disarmament during the time of Metal Gear 2 
	was omitted from the Previous Operations section of Metal Gear Solid.`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("positive test", func(t *testing.T) {
		expected := []string{
			"the",      //21
			"of",       //12
			"to",       //9
			"a",        //7
			"and",      //7
			"was",      //7
			"by",       //6
			"nuclear",  //6
			"Zanzibar", //5
			"Gear",     //3
		}
		require.Equal(t, expected, Top10(textEnglish))
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.Equal(t, expected, Top10(text))
		}
	})
}
