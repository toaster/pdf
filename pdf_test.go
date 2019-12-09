package pdf_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/toaster/pdf"
)

func TestPdf_ExtractText(t *testing.T) {
	tests := map[string]struct {
		src       string
		wantPages int
		want      []string
	}{
		"multiple content streams per page": {
			src:       "testdata/multi_stream_content.pdf",
			wantPages: 6,
			want: []string{
				"Are You Asking Your CMSVendor the Right Questions?scrivito.com",
				"scrivito.com2How does the CMS help me to improve website projects?How does the CMS support team collaboration in the project and later?How does the CMS a\x1fect my business case?Is the UI easy to understand and use without training?How does the CMS deal with customer experience?Where is the limit in visitors/users/transactions per month for the CMS?Is the CMS limiting my ideas for the web project?What does the CMS do to support SEO?What does the CMS o\x1fer to ensure fast page loading times?   Is there a need to involve the IT administration?   1.2.3.4.5.6.7.8.9.10.Choosing the right Content Management System (CMS) for your web project at first could appear to bea momentous task. However, if you ask the right questions before you get too deep, your whole project will come together more easily than you can imagine and you\u0092ll be well on your way to a successful launch.Here are ten things you can ask of the CMS vendor to help you choose the right modern Web CMS.",
				"3How does the CMShelp me to improve website projects?When it comes to selecting a CMS in today's modern environment, you really have three main choices to pickbetween. First, you can choose to go with a basic but functional CMS like WordPress. Secondly, you can chooseto utilize a more advanced CMS framework on the other end of the complexity spectrum, with Joomla being justone of many examples. Finally, you may choose to implement the needed CMS functionality on your own accord.In order to be sure that you're making the right decision, you need to look out for a few key characteristics.The CMS you choose needs to yield instant results through the entirety of the project lifecycle, from the first pitchto its eventual launch. Any high quality CMS deployment should also have on-demand availability. You shouldn'tneed to install any so\x1fware, manage any upgrades on your own, look for available plug-ins and more.You certainly shouldn't need to purchase any additional hardware for your business to get the CMS o\x1e the ground.You should be able to load up your browser, sign in and start working.The CMS vendor that you choose should also be able to display the flexibility that you will need moving forward.It should be easy to integrate and even easier to extend this functionality on an as-needed basis. Your CMS shouldalso have a high degree of useful, prefabricated templates and widgets and support the integration of popularweb frameworks like Bootstrap, for example.Finally, your CMS solution should be able to grow and evolve as your business does the same. If you needto be able to reuse certain components from other projects, be they third party elements or your own, your CMSshould allow you to do that. Most importantly, your CMS should require minimal operating e\x1eort through allstages of the project to truly allow you to create the high quality content that you're a\x1fer.scrivito.comHow does the CMS support teamcollaboration in the project and later?Nowadays, you rarely can limit access to a system to a fixed set of physical locations. You need to be browser-based,permitting access from anywhere. How does the CMS support this while supporting the division of work betweenan arbitrary number of editors, designers, and programmers? Is there a concept of an editorial team where youhave workspaces and previews before being released to the masses? Does the CMS detect conflicts deriving fromconcurrent editing and how does the CMS support the resolution of such conflicts? 1.2.",
				"4How does the CMSa\x1fect my business case?When choosing a CMS, you typically want something with a low cost to get started. You may want to create a proofof concept project to help get buy-in from the company, but you want a low cost and minimal investment for thisearly stage project. Is the system pay-as-you-go or does it require dedicated hardware and server so\x1fware?Can you get started without having to add new IT sta\x1e or fixed data center fees? The higher the up front cost,the higher the risk in choosing a CMS. Everything might look great on paper, but if you can\u0092t kick the tires a little,you\u0092ll never know for sure if the choice is right, if the startup costs are exorbitant, and you need upper managementapproval to take the costly risk.scrivito.comIs the UI easy to understandand use without training?Some personnel balk at the concept of adding a CMS to their everyday project workflow. If they have to take aweek of training before they get started, not only are users behind a week but they quickly learn why they neededthat week of training. The system may just be too complicated and requires a long ramp up time. A Drag and DropWYSIWYG system will get a team up and running much more quickly.How does the CMS deal withcustomer experience?A website run by a CMS requires some extra work to ensure users can find and use the content, but potentiallymore importantly, the search engines can find and read the content. Is there internationalization, localization,and personalization support available? What about SEO support? How much of the meta tags can be configuredand how much of the support files can be automated? This also deals with responsive design if your users cancome through from multiple platforms and devices. You don\u0092t want to lock out mobile users and deal with themin version two of the system.3.4.5.",
				"5Where is the limit in visitors/users/transactionsper month for the CMS?Is the CMS size constrained? If your project takes o\x1f in popularity, do you have to change CMSs because you\u0092veoutgrown infrastructure limitations? Are there upper limits for tra\x1fic? Obviously, something like Twitter andPinterest require some level of customization given their current sizes. But, in their early days, were they flexibleenough to handle their instant growth patterns? Is scalability an option? As rich content assets like videos orhigh-res images become more and more popular, is the CMS limited in the number of content objects or in thesize of these objects?scrivito.comIs the CMS limiting my ideasfor the web project?Does the CMS restrict you from doing what you want with the web site or are you constrained to do everythingwithin the CMS? Is the CMS an open framework, o\x1fering the building blocks for a successful system or a closedenvironment? You don\u0092t want the CMS to get in your way as you expand beyond your initial needs.What does the CMS doto support SEO?Being found by search engines is important for most websites. A CMS might have a feature checkbox for SEOsupport, but what exactly does that SEO support involve? Does it have support for a full set of features for on-pageoptimizations? Is the URL structure SEO friendly? Is there XML sitemap creation support? Does it avoid Flash?Are there tools for editors? The more help you can get from the CMS, the better optimized a site could be.You don\u0092t need to use all the features immediately, but you want them there.6.7.8.",
				"6What does the CMS o\x1ferto ensure fast page loading times?Loading times are incredibly important in today's modern environment. Studies have shown that users willtypically only spend three seconds or less waiting for a page to load, at which point they will almost certainly gosomewhere else to find the information that they're looking for. As a result, load times should always be one ofyour primary concerns and the CMS that you choose needs to reflect that idea.When searching for a CMS, you should be looking for one that not only o\x1fers automatic load balancing, cachingand fast loading times via a CDN, but other advanced features at well. Instead of shunning multimedia elementslike rich media or video in fear of what they might do to your load times, you should be instead searching fora CMS deployment that will allow you to incorporate these elements and still have the lightning fast page thatyou can depend on.Your CMS should also include support for an elastic infrastructure, meaning that it should be able to scale content,users and workspaces as the needs of your site continue to change. Instead of planning for the infrastructurethat you think you're going to have tomorrow, elastic infrastructure and CMS in general can help you accountfor the specifically needed demand of the present. This will ultimately create a much more functional,higher quality site as a result.Is there a need to involvethe IT administration?Lastly, how much need is there to involve the IT administration? Beyond just hiring sta\x1f who understands theintricacies of the CMS, you may need to have to add more security and backup requirements. With a so\x1eware asa service (SaaS) approach, there are no so\x1eware packages that need to be installed and maintained, and noadditional IT infrastructure is needed. The number of questions you should ask before choosing a modernWeb CMS is endless. These ten should provide a good basis for you making the right decision for your company.9.10.Sign up for a free trial for Scrivito at https://scrivito.com Scrivito is proudly made by Infopark in Berlin, Germanyscrivito.com",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			data, err := ioutil.ReadFile(tt.src)
			require.NoError(t, err)
			bufReader := bytes.NewReader(data)
			pdfReader, err := pdf.NewReader(bufReader, bufReader.Size())
			if assert.NoError(t, err) {
				totalPage, err := pdfReader.NumPage()
				if assert.NoError(t, err) {
					assert.Equal(t, tt.wantPages, totalPage)
					var texts []string
					for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
						page, err := pdfReader.Page(pageIndex)
						if assert.NoError(t, err) {
							assert.NotNil(t, page.V)
							content, err := page.Content()
							if assert.NoError(t, err) {
								var text string
								for _, t := range content.Text {
									text += t.S
								}
								texts = append(texts, text)
							}
						}
					}
					assert.Equal(t, tt.want, texts)
				}
			}
		})
	}
}
