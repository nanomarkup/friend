package main

import "os"

func getFeeds() []*feed {
	wd, _ := os.Getwd()
	feeds := []*feed{}
	// 1 Загальні положення
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802081", wd + "/data/6802081.nam"})
	// 2.1 Повноваження та персонал: призначення, звільнення та інші ситуації
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802084", wd + "/data/6802084.nam"})
	// 2.2 Влада та персонал: курси, протистояння та змагання
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802085", wd + "/data/6802085.nam"})
	// 2.3 Повноваження та персонал: інші
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802086", wd + "/data/6802086.nam"})
	// 3 Адміністративний договір
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802087", wd + "/data/6802087.nam"})
	// 4.1 Економіка та фінанси: дії в бюджетних питаннях
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802089", wd + "/data/6802089.nam"})
	// 4.2 Економіка та фінанси: дії у фіскальних справах
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802090", wd + "/data/6802090.nam"})
	// 4.3 Економіка та фінанси: дії щодо соціального забезпечення
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802091", wd + "/data/6802091.nam"})
	// 4.4 Економіка та фінанси: інше
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802092", wd + "/data/6802092.nam"})
	// 5 Примусове відчуження
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802094", wd + "/data/6802094.nam"})
	// 6 Субсидії та допомоги
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802095", wd + "/data/6802095.nam"})
	// 7.1 Інші оголошення: Містобудування
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802097", wd + "/data/6802097.nam"})
	// 7.2 Інші оголошення: навколишнє середовище та енергетика
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802098", wd + "/data/6802098.nam"})
	// 7.3 Інші оголошення: статути та колективні договори
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802099", wd + "/data/6802099.nam"})
	// 7.4 Інші оголошення: фізичні особи
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802100", wd + "/data/6802100.nam"})
	// 7.5 Інші оголошення: різні
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802301", wd + "/data/6802301.nam"})
	// 8.1 Судові процедури: Аукціони
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/7479572", wd + "/data/7479572.nam"})
	// 8.2 Судові процедури: інші оголошення
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802303", wd + "/data/6802303.nam"})
	// 9 Вибори
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/7293890", wd + "/data/7293890.nam"})
	// Міністерство президента, юстиції, безпеки та адміністративного спрощення
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/27148/size20", wd + "/data/27148.nam", 4})
	// Міністерство розвитку, територіального планування та екології
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/27147/size20", wd + "/data/27147.nam", 4})
	// Міністерство економіки, фінансів та європейських фондів
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/27141/size20", wd + "/data/27141.nam", 4})
	// Міністерство освіти, професійної підготовки та університетів
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/27139/size20", wd + "/data/27139.nam", 4})
	// Міністерство культури, туризму та спорту
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/27146/size20", wd + "/data/27146.nam", 4})
	// Міністерство сільського розвитку, тваринництва, рибальства та продовольства
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/27140/size20", wd + "/data/27140.nam", 4})
	// Міністерство промисловості, зайнятості, інновацій та торгівлі
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/27145/size20", wd + "/data/27145.nam", 4})
	// Консультації щодо здоров'я
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/27149/size20", wd + "/data/27149.nam", 4})
	// Департамент соціального залучення, молоді, сім'ї та рівноправності
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/7720252/size20", wd + "/data/7720252.nam", 4})
	// Остання допомога та субсидії
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16401/eid4008203/size50", wd + "/data/4008203.nam", 4})
	// Останні стипендії
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16401/eid4008202/size50", wd + "/data/4008202.nam", 4})
	// Останні нагороди
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16401/eid4008221/size50", wd + "/data/4008221.nam", 4})
	// Останні процедури
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16401/eid4008246/size50", wd + "/data/4008246.nam", 4})
	// Останні прес-релізи
	// Disabled because it duplicates 27139.nam
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/size50", wd + "/data/4008216.nam", 4})
	// Новини пропозиції про державну роботу
	feeds = append(feeds, &feed{"https://empleopublico.cantabria.es/o/GOBIERNO/feed/group16475/inscom_liferay_journal_content_web_portlet_JournalContentPortlet_INSTANCE_6Cx0YFAD8ZVK/size50", wd + "/data/6Cx0YFAD8ZVK.nam", 4})
	return feeds
}
