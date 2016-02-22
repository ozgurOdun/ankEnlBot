package tgHandler

import (
	"bytes"
	"dbOps"
	"fmt"
	godotenv "github.com/joho/godotenv"
	"github.com/rockneurotiko/go-tgbot"
	"os"
	"strconv"
)

var availableCommands = map[string]string{
	"/gorevekle":        "Listeye yeni gorev ekle",
	"/gorevtamamlandi":  "Gorevi yapildi olarak isaretle",
	"/goreviptal":       "Gorevi listeden silme",
	"/gorevlerilistele": "Acik gorevleri listele",
	"/gorevklavyesi":    "Bota ozel bir klavye acar",
	"/start":            "Botu baslatir",
	"/yardim":           "Bot komutlarini listeler", ///done
	"/senkimsin":        "Info",                     ///done
}

func BotInit() *tgbot.TgBot {
	godotenv.Load("token.env")
	///token.env contains the token for telegram bot like:
	/// TELEGRAM_KEY=token

	token := os.Getenv("TELEGRAM_KEY")
	bot := tgbot.NewTgBot(token)
	bot.DefaultDisableWebpagePreview(true) // Disable all link preview by default
	bot.DefaultOneTimeKeyboard(true)       // Enable one time keyboard by default
	bot.DefaultSelective(true)             // Use Seletive by default

	bot.
		SimpleCommandFn(`yardim`, buildHelpMessage).
		SimpleCommandFn(`klavyeyisakla`, hideKeyboard).
		SimpleCommandFn(`senkimsin`, showAbout).
		CommandFn(`gorevekle (.+)`, addTaskHandler).
		CommandFn(`gorevtamamlandi (\d)`, markCompleteHandler).
		CommandFn(`goreviptal (\d)`, deleteTaskHandler).
		MultiCommandFn([]string{`gorevlerilistele`, `gorevlerilistele(\s)(.+)?`}, listTasksHandler).
		AnyMsgFn(allMsgHand)
	return bot
}

func allMsgHand(bot tgbot.TgBot, msg tgbot.Message) {
	// comment this to unsee it :)
	fmt.Printf("Received message: %+v\n", msg)
	exists := dbOps.CheckUser(*msg.From.Username)
	if exists == false {
		if msg.From.LastName == nil {
			msg.From.LastName = new(string)
		}
		if &msg.From.FirstName == nil {
			msg.From.FirstName = "null"
		}
		if msg.From.Username == nil {
			msg.From.Username = new(string)
			*msg.From.Username = strconv.Itoa(msg.From.ID)
			fmt.Printf("test\n")
		}
		fmt.Printf(*msg.From.LastName)
		added := dbOps.AddNewUser(msg.From.ID, *msg.From.Username, msg.From.FirstName, *msg.From.LastName)
		if added == false {
			fmt.Printf("Adding new user failed\n")
		} else {
			fmt.Printf("New user added:%s\n", *msg.From.Username)
			message := "Merhaba " + *msg.From.Username + " Seni de listeme ekledim."
			bot.Answer(msg).Text(message).End()
		}
	} else {
		fmt.Printf("User %s already exists\n", *msg.From.Username)
	}

	//query database at user tables if user exists msg.From.Username
	//if not add a new row to users table and give user an id
}

func addTaskHandler(bot tgbot.TgBot, msg tgbot.Message, vals []string, kvals map[string]string) *string {
	msgtext := ""
	if len(vals) > 1 {
		msgtext = vals[1]
		fmt.Printf(msgtext + "\n")
	}
	user := dbOps.GetUserByName(*msg.From.Username)
	if user != nil {
		ret := dbOps.AddNewTask(user.Uid, msgtext)
		if ret > 0 {
			bot.Answer(msg).Text("Yeni görev başarı ile eklendi").End()
		}
	}
	return nil
}
func markCompleteHandler(bot tgbot.TgBot, msg tgbot.Message, vals []string, kvals map[string]string) *string {
	if len(vals) > 1 {
		if id, err := strconv.Atoi(vals[1]); err == nil {
			dbOps.MarkTaskAsDone(id)
			message := strconv.Itoa(id) + " numaralı görev tamamlandı ve bir daha listede görünmeyecek."
			bot.Answer(msg).Text(message).End()
		} else {
			bot.Answer(msg).Text("/gorevtamamalandi numara şeklinde kullanabilirsin.").End()
		}
	}
	return nil
}

func deleteTaskHandler(bot tgbot.TgBot, msg tgbot.Message, vals []string, kvals map[string]string) *string {
	if len(vals) > 1 {
		if id, err := strconv.Atoi(vals[1]); err == nil {
			if dbOps.DeleteTask(id) {
				message := strconv.Itoa(id) + " numaralı görev iptal edildi"
				bot.Answer(msg).Text(message).End()
			}
		} else {
			bot.Answer(msg).Text("/goreviptal numara şeklinde kullanabilirsin.").End()
		}
	}
	return nil
}

func listTasksHandler(bot tgbot.TgBot, msg tgbot.Message, vals []string, kvals map[string]string) *string {
	var tasks []*dbOps.Task
	var num int64
	var index int64
	var message string
	if len(vals) > 2 {
		exists := dbOps.CheckUser(vals[2])
		if exists {
			user := dbOps.GetUserByName(vals[2])
			bot.Answer(msg).Text(user.UserName + " isimli ajana ait açık görevleri hemen listeliyorum.").End()
			tasks, num = dbOps.ListUndoneTasksByUser(user.UserName)
			if num == 0 {
				bot.Answer(msg).Text("Bu ajana ait hiç görev bulamadım. /gorevekle bişiyler komutu ile yeni görev" +
					" ekleyebilirsin").End()
				return nil
			}
			user = nil
			user = dbOps.GetUserById(tasks[index].Ownerid)
			for index = 0; index < num; index++ {
				message = strconv.Itoa(tasks[index].Uid) + "-" + tasks[index].Task + " " +
					tasks[index].CreationTime.Format("01.02.2006 15:04") + " tarihinde " + user.UserName +
					" tarafından oluşturuldu."
				bot.Answer(msg).Text(message).End()
			}
			bot.Answer(msg).Text("Son.").End()
			bot.Answer(msg).Text("/gorevtamamlandi #görevnumarası ya da /goreviptal #görevnumarası komutlarını kullanarak " +
				"bişiler yapabilirsin").End()
		} else {
			bot.Answer(msg).Text("Afedersin böyle bir kullanıcı bulamadım. Lütfen ismi kontrol eder misin?").End()
		}
	} else {
		bot.Answer(msg).Text("Açık görevleri hemen listeliyorum.").End()
		tasks, num = dbOps.ListUndoneTasks()
		if num == 0 {
			bot.Answer(msg).Text("Hiç görev bulamadım. /gorevekle bişiyler komutu ile yeni görev" +
				" ekleyebilirsin").End()
			return nil
		}
		user := dbOps.GetUserById(tasks[index].Ownerid)
		for index = 0; index < num; index++ {
			message = strconv.Itoa(tasks[index].Uid) + "-" + tasks[index].Task + " " +
				tasks[index].CreationTime.Format("01.02.2006 15:04") + " tarihinde " + user.UserName +
				" tarafından oluşturuldu."
			bot.Answer(msg).Text(message).End()
		}
		bot.Answer(msg).Text("Son.").End()
		bot.Answer(msg).Text("/gorevtamamlandi #görevnumarası ya da /goreviptal #görevnumarası komutlarını kullanarak " +
			"bişiler yapabilirsin").End()
		bot.Answer(msg).Text("Ayrıca /gorevlerilistele username yazarak tek bir kullanıcının açtığı görevleri " +
			"görüntüleyebilirsin").End()
	}
	return nil
}

func showAbout(bot tgbot.TgBot, msg tgbot.Message, text string) *string {
	bot.Answer(msg).Text("Merhaba, ben Anka-RA Enlightened için ozgurOdun tarafından yazılmış bir görev düzenleme " +
		"botuyum. /yardim a tıklayarak neler yapabileceğimi öğrenebilirsin.").End()
	return nil
}

/**
* Hides the keyboard and reports "done"
 */
func hideKeyboard(bot tgbot.TgBot, msg tgbot.Message, text string) *string {
	rkm := tgbot.ReplyKeyboardHide{HideKeyboard: true, Selective: false}
	bot.Answer(msg).Text("Tamam, özel klavyeyi gizliyorum!").KeyboardHide(rkm).End()
	return nil
}

/**
* builds help message displayed when /help command is called
 */
func buildHelpMessage(bot tgbot.TgBot, msg tgbot.Message, arg string) *string {
	var buffer bytes.Buffer
	for cmd, htext := range availableCommands {

		str := fmt.Sprintf("%s - %s\n", cmd, htext)
		buffer.WriteString(str)
	}
	s := buffer.String()
	return &s
}
